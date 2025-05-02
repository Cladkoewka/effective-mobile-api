package main

import (
	"flag"
	"fmt"
	"os"

	_ "github.com/Cladkoewka/effective-mobile-api/docs"
	"github.com/Cladkoewka/effective-mobile-api/internal/api/enrichment"
	"github.com/Cladkoewka/effective-mobile-api/internal/config"
	"github.com/Cladkoewka/effective-mobile-api/internal/handler"
	"github.com/Cladkoewka/effective-mobile-api/internal/logger"
	"github.com/Cladkoewka/effective-mobile-api/internal/repository"
	"github.com/Cladkoewka/effective-mobile-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	logger.InitLogger()

	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	logger.Log.Info("Starting application")

	cfg := initConfig()
	db := initDatabase(cfg)

	if *migrateFlag {
		logger.Log.Info("Running database migration")
		if err := runMigration(db, "migrations/001_init.sql"); err != nil {
			logger.Log.Error("Migration failed", "error", err)
		}
		logger.Log.Info("Migration completed successfully")
	}

	runServer(db, cfg)
}

func initConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Log.Error("Error loading config", "error", err)
	}
	return cfg
}

func initDatabase(cfg *config.Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres",
		"postgres://"+cfg.DBUser+":"+cfg.DBPassword+"@"+cfg.DBHost+":"+cfg.DBPort+"/"+cfg.DBName+"?sslmode=disable")
	if err != nil {
		logger.Log.Error("Error connecting to the database", "error", err)
	}
	return db
}

func runMigration(db *sqlx.DB, filePath string) error {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading migration file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}

	logger.Log.Info("Migration applied from file", "file", filePath)
	return nil
}

func runServer(db *sqlx.DB, cfg *config.Config) {
	personRepo := repository.NewPersonRepository(db)
	enrichmentClient := enrichment.NewHTTPEnrichmentClient()
	personService := service.NewPersonService(personRepo, enrichmentClient)
	personHandler := handler.NewPersonHandler(personService)

	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	personHandler.RegisterRoutes(api)

	logger.Log.Info("Server starting", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Log.Error("Error starting the server", "error", err)
	}
}
