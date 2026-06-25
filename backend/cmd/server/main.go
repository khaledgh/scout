package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"kashfi/internal/config"
	"kashfi/internal/db"
	"kashfi/internal/jobs"
	appMiddleware "kashfi/internal/middleware"
	"kashfi/internal/ws"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}

	if cfg.App.Env == "development" {
		if err := db.AutoMigrate(database); err != nil {
			log.Fatalf("db migrate: %v", err)
		}
		log.Println("db migration complete")
	}

	hub := ws.NewHub()
	go hub.Run()

	scheduler := jobs.New(database)
	scheduler.Start()
	defer scheduler.Stop()

	e := echo.New()
	e.HideBanner = true

	// Global middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(appMiddleware.CORS(cfg.CORS.AllowedOrigins))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		HSTSMaxAge:            31536000,
		ContentSecurityPolicy: "",
	}))

	// Static uploads
	e.Static(cfg.Upload.PublicPath, cfg.Upload.Dir)

	// Serve frontend static assets (JS, CSS, icons, etc.)
	e.Static("/assets", "./public/assets")
	e.Static("/icons",  "./public/icons")

	// API group
	api := e.Group("/api/v1")

	// Register routes
	registerRoutes(api, cfg, database, hub)

	// SPA fallback — all non-API, non-asset routes return index.html
	e.GET("/*", func(c echo.Context) error {
		return c.File("./public/index.html")
	})

	// Start server with graceful shutdown
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on port %s (env=%s)", cfg.App.Port, cfg.App.Env)
		if err := e.StartServer(srv); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	log.Println("server stopped")
}
