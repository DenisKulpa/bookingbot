package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/lib/pq"
)

func New(dataSourceName, dbName string) (*sql.DB, error) {
	log.Printf("db: connecting to postgres")

	// Пробуем подключиться к целевой базе
	database, err := tryConnect(dataSourceName)
	if err != nil && isDBMissing(err) {
		// Базы нет — создаём через подключение к postgres
		log.Printf("db: database %s not found, creating...", dbName)

		sysDSN := replaceDBName(dataSourceName, "postgres")
		sysDB, cerr := sql.Open("postgres", sysDSN)
		if cerr != nil {
			return nil, fmt.Errorf("db.New connect to postgres: %w", cerr)
		}
		defer sysDB.Close()

		if _, cerr := sysDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s"`, dbName)); cerr != nil {
			return nil, fmt.Errorf("db.New create database: %w", cerr)
		}
		log.Printf("db: database %s created", dbName)

		// Переподключаемся к созданной базе
		database, err = tryConnect(dataSourceName)
	}
	if err != nil {
		return nil, err
	}

	database.SetMaxOpenConns(25)
	database.SetMaxIdleConns(5)

	if err := runMigrations(database); err != nil {
		return nil, fmt.Errorf("db migrations: %w", err)
	}

	return database, nil
}

func tryConnect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("db.New: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}

// isDBMissing проверяет, что ошибка — «база данных не существует» (SQLSTATE 3D000).
func isDBMissing(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "3D000"
	}
	return false
}

// replaceDBName заменяет dbname в key=value DSN
func replaceDBName(dsn, newName string) string {
	parts := strings.Fields(dsn)
	for i, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			parts[i] = "dbname=" + newName
			return strings.Join(parts, " ")
		}
	}
	return dsn + " dbname=" + newName
}

func runMigrations(db *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(filename), "..", "..", "migrations")

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    TEXT PRIMARY KEY,
		applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	)`)
	if err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	if err != nil {
		return fmt.Errorf("glob migrations: %w", err)
	}
	sort.Strings(files)

	for _, file := range files {
		version := filepath.Base(file)

		var count int
		if err := db.QueryRow(`SELECT COUNT(*) FROM schema_migrations WHERE version = $1`, version).Scan(&count); err != nil {
			return fmt.Errorf("check migration %s: %w", version, err)
		}
		if count > 0 {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", version, err)
		}

		tx, err := db.Begin()
		if err != nil {
			return fmt.Errorf("begin transaction for %s: %w", version, err)
		}

		// Выполняем каждый оператор отдельно (pq не поддерживает multi-statement Exec)
		for _, stmt := range splitSQL(string(content)) {
			if _, err := tx.Exec(stmt); err != nil {
				_ = tx.Rollback()
				return fmt.Errorf("apply migration %s: %w", version, err)
			}
		}

		if _, err := tx.Exec(`INSERT INTO schema_migrations (version) VALUES ($1)`, version); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("record migration %s: %w", version, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("commit migration %s: %w", version, err)
		}

		log.Printf("applied migration: %s", version)
	}

	return nil
}

// RunSeed выполняет seeds/seed.sql — идемпотентно (все запросы ON CONFLICT DO NOTHING).
func RunSeed(database *sql.DB) error {
	_, filename, _, _ := runtime.Caller(0)
	seedFile := filepath.Join(filepath.Dir(filename), "..", "..", "seeds", "seed.sql")

	content, err := os.ReadFile(seedFile)
	if err != nil {
		return fmt.Errorf("read seed file: %w", err)
	}

	for _, stmt := range splitSQL(string(content)) {
		if _, err := database.Exec(stmt); err != nil {
			return fmt.Errorf("run seed: %w", err)
		}
	}

	log.Printf("db: seed applied from %s", filepath.Base(seedFile))
	return nil
}

// splitSQL разбивает SQL-текст на отдельные операторы по ; (игнорируя пустые строки и комментарии).
func splitSQL(sql string) []string {
	var stmts []string
	var current []byte
	inString := false
	inDollar := false

	for i := 0; i < len(sql); i++ {
		c := sql[i]

		// Отслеживаем строки и доллар-квотинг
		if c == '\'' && !inDollar {
			inString = !inString
		}
		if i+1 < len(sql) && sql[i] == '$' && sql[i+1] == '$' && !inString {
			inDollar = !inDollar
			current = append(current, c)
			i++
			current = append(current, sql[i])
			continue
		}

		if c == ';' && !inString && !inDollar {
			stmt := trimSQL(string(current))
			if stmt != "" {
				stmts = append(stmts, stmt)
			}
			current = nil
			continue
		}

		current = append(current, c)
	}

	stmt := trimSQL(string(current))
	if stmt != "" {
		stmts = append(stmts, stmt)
	}

	return stmts
}

func trimSQL(s string) string {
	s = strings.TrimSpace(s)
	// Убираем строки-комментарии (-- ...)
	var lines []string
	for _, line := range splitLines(s) {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && !strings.HasPrefix(trimmed, "--") {
			lines = append(lines, line)
		}
	}
	return strings.TrimSpace(strings.Join(lines, "\n"))
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
