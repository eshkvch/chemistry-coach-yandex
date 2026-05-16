// @title           Alice Chemistry Coach API
// @version         1.0
// @description     18+ AI-тренажёр романтической коммуникации (MVP backend)
// @BasePath        /api/v1
// @host            localhost:8080
// @schemes         http
package main

import (
	"fmt"
	"log"
	"path/filepath"

	"chemistry-coach/internal/config"
	httpx "chemistry-coach/internal/delivery/http"
	"chemistry-coach/internal/delivery/http/handler"
	"chemistry-coach/internal/domain"
	"chemistry-coach/internal/infrastructure/postgres"
	"chemistry-coach/internal/infrastructure/yandex"
	"chemistry-coach/internal/usecase"

	_ "chemistry-coach/docs"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := postgres.Connect(cfg.DatabaseURL, cfg.IsDevelopment())
	if err != nil {
		log.Fatal(err)
	}
	migrationsDir := filepath.Join(".", "migrations")
	if err := postgres.RunMigrations(db, migrationsDir); err != nil {
		log.Fatalf("migration: %v", err)
	}
	if err := postgres.Ping(db); err != nil {
		log.Fatalf("db ping: %v", err)
	}

	var llm domain.LLMService
	if cfg.Yandex.APIKey != "" && cfg.Yandex.FolderID != "" {
		llm = yandex.NewClient(cfg.Yandex)
		log.Println("using Yandex AI Studio")
	} else {
		llm = yandex.NewMockClient()
		log.Println("Yandex AI not configured — using mock LLM")
	}

	userRepo := postgres.NewUserRepo(db)
	sessionRepo := postgres.NewSessionRepo(db)
	messageRepo := postgres.NewMessageRepo(db)
	debriefRepo := postgres.NewDebriefRepo(db)

	authUC := usecase.NewAuthUseCase(userRepo)
	profileUC := usecase.NewProfileUseCase(userRepo, sessionRepo, messageRepo)
	catalogUC := usecase.NewCatalogUseCase()
	sessionUC := usecase.NewSessionUseCase(sessionRepo, messageRepo, debriefRepo, llm)

	router := httpx.NewRouter(httpx.Handlers{
		Auth:    handler.NewAuthHandler(authUC),
		Profile: handler.NewProfileHandler(profileUC),
		Catalog: handler.NewCatalogHandler(catalogUC),
		Session: handler.NewSessionHandler(sessionUC),
	}, cfg.IsDevelopment())

	addr := fmt.Sprintf(":%d", config.ParsePort(cfg.Port))
	log.Printf("listening on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
