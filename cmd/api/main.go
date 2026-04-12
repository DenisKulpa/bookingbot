package main

import (
	"log"
	"net/http"

	"github.com/DenisKulpa/bookingbot/internal/config"
	"github.com/DenisKulpa/bookingbot/internal/db"
	"github.com/DenisKulpa/bookingbot/internal/handler"
	"github.com/DenisKulpa/bookingbot/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.New(cfg.SQLitePath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer database.Close()

	zoneRepo := repository.NewZoneRepository(database)
	zoneHandler := handler.NewZoneHandler(zoneRepo)

	apartmentRepo := repository.NewApartmentRepository(database)
	apartmentHandler := handler.NewApartmentHandler(apartmentRepo)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Get("/districts", zoneHandler.GetDistricts)
		r.Get("/districts/{id}", zoneHandler.GetDistrictDetail)
		r.Get("/districts/{id}/apartments", apartmentHandler.GetApartments)
	})

	log.Printf("Server started on :%s", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, r); err != nil {
		log.Fatalf("server: %v", err)
	}
}
