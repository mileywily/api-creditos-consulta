package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"api-creditos-consulta/internal/app/api/handler"
	"api-creditos-consulta/internal/core/usecases"
	"api-creditos-consulta/internal/infra/client"
)

func main() {
	// Setup JSON Logger (Replicando la estructura de Logback JSON Layout del legado)
	replaceAttr := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey {
			a.Key = "timestamp"
			a.Value = slog.StringValue(a.Value.Time().Format("2006-01-02T15:04:05.000Z"))
		}
		if a.Key == slog.MessageKey {
			a.Key = "message"
		}
		if a.Key == slog.LevelKey {
			// keep "level"
		}
		return a
	}
	
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		ReplaceAttr: replaceAttr,
		Level:       slog.LevelDebug, // Legado usa Debug en logging.level
	}))
	slog.SetDefault(logger)

	// Validate critical environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	finnflowURL := os.Getenv("FINNFLOW_URL")
	if finnflowURL == "" {
		slog.Error("CRITICAL: FINNFLOW_URL environment variable is missing")
		os.Exit(1)
	}
	finnflowKey := os.Getenv("FINNFLOW_KEY")
	if finnflowKey == "" {
		slog.Error("CRITICAL: FINNFLOW_KEY environment variable is missing")
		os.Exit(1)
	}
	finnflowSecret := os.Getenv("FINNFLOW_SECRET")
	if finnflowSecret == "" {
		slog.Error("CRITICAL: FINNFLOW_SECRET environment variable is missing")
		os.Exit(1)
	}

	slog.Info("Starting api-creditos-consulta", "port", port, "finnflow_url", finnflowURL)

	// Start APM Tracer
	tracer.Start(
		tracer.WithService("api-creditos-consulta"),
		tracer.WithEnv(os.Getenv("ENV")),
	)
	defer tracer.Stop()

	// Initialize Gin
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = v.RegisterValidation("customRutPattern", func(fl validator.FieldLevel) bool {
			matched, _ := regexp.MatchString(`^[0-9]{1,2}\.[0-9]{3}\.[0-9]{3}-[0-9kK]$`, fl.Field().String())
			// También soportamos sin puntos para mayor flexibilidad, o mantenemos estricto como el regex
			if !matched {
				matched, _ = regexp.MatchString(`^[0-9]+-[0-9kK]$`, fl.Field().String())
			}
			return matched
		})
	}
	
	// Initialize dependencies
	finnflowClient := client.NewFinnFlowClient(finnflowURL, finnflowKey, finnflowSecret)
	creditosUseCase := usecases.NewCreditosUseCase(finnflowClient)
	creditosHandler := handler.NewCreditosHandler(creditosUseCase)

	// Middlewares
	router.Use(gintrace.Middleware("api-creditos-consulta"))
	router.Use(gin.Recovery())

	// Healthcheck
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "up"})
	})

	// Endpoints
	api := router.Group("/api")
	{
		creditos := api.Group("/creditos")
		{
			creditos.POST("/consulta/", creditosHandler.ConsultarCreditos)
		}
	}

	// Server setup
	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", "error", err)
	}
	slog.Info("Server exiting")
}
