package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/PashakArt/file-server/internal/config"
	"github.com/PashakArt/file-server/internal/db/redis"
	"github.com/PashakArt/file-server/internal/db/repository"
	"github.com/PashakArt/file-server/internal/service"
	"github.com/PashakArt/file-server/internal/transport/http/handler"
	"github.com/PashakArt/file-server/internal/transport/http/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbPool, err := pgxpool.New(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer dbPool.Close()

	rdb, err := redis.NewRedisClient(*cfg)
	if err != nil {
		log.Fatalf("redis error: %v", err)
	}
	defer rdb.Close()

	userRepo := repository.NewUserRepository(dbPool)
	docRepo := repository.NewDocRepository(dbPool)

	authService := service.NewAuthService(userRepo, rdb, cfg.TokenTTL)
	docService := service.NewDocService(cfg, docRepo, userRepo)

	authMiddleware := middleware.AuthMiddleware(authService)

	authHandler := handler.NewAuthHandler(authService, cfg.AdminToken)
	docHandler := handler.NewDocsHandler(docService, cfg)

	mux := http.NewServeMux()
	authHandler.RegisterRoutes(mux)
	docHandler.RegisterRoutes(mux, authMiddleware)

	addr := fmt.Sprintf(":%s", cfg.HttpPort)
	log.Printf("Server running on http://localhost%s", addr)

	err = http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
