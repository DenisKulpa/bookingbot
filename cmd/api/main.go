package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/DenisKulpa/bookingbot/internal/config"
	"github.com/DenisKulpa/bookingbot/internal/db"
	"github.com/DenisKulpa/bookingbot/internal/handler"
	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/DenisKulpa/bookingbot/internal/telegram"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	if cfg.Seed {
		if err := db.RunSeed(database); err != nil {
			log.Fatalf("seed: %v", err)
		}
	}

	zoneRepo := repository.NewZoneRepository(database)
	zoneHandler := handler.NewZoneHandler(zoneRepo)

	apartmentRepo := repository.NewApartmentRepository(database)
	apartmentHandler := handler.NewApartmentHandler(apartmentRepo)

	// ── Telegram bot ──────────────────────────────────────────────────────────
	tgClient, err := telegram.New(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("telegram: %v", err)
	}
	bot := telegram.NewBot(tgClient, zoneRepo, apartmentRepo)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go bot.Run(ctx)

	// ── HTTP API ──────────────────────────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/districts", zoneHandler.GetDistricts)
		r.Get("/districts/{id}", zoneHandler.GetDistrictDetail)
		r.Get("/districts/{id}/apartments", apartmentHandler.GetApartments)
		r.Get("/apartments/{id}", apartmentHandler.GetApartmentDetail)
	})

	log.Printf("Server started on :%s", cfg.ServerPort)
	srv := &http.Server{Addr: ":" + cfg.ServerPort, Handler: r}

	go func() {
		<-ctx.Done()
		log.Println("shutting down HTTP server...")
		_ = srv.Shutdown(context.Background())
	}()

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server: %v", err)
	}
}
