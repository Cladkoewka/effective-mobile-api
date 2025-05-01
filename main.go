package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"flag"

	_ "github.com/Cladkoewka/effective-mobile-api/docs"
	"github.com/Cladkoewka/effective-mobile-api/internal/api/enrichment"
	"github.com/Cladkoewka/effective-mobile-api/internal/config"
	"github.com/Cladkoewka/effective-mobile-api/internal/handler"
	"github.com/Cladkoewka/effective-mobile-api/internal/repository"
	"github.com/Cladkoewka/effective-mobile-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // PostgreSQL driver
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migrations")
	flag.Parse()

	cfg := initConfig()

	db := initDatabase(cfg)

	if *migrateFlag {
		if err := runMigration(db, "migrations/001_init.sql"); err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully")
	}

	runServer(db, cfg)
}

func initConfig() *config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	return cfg
}

func initDatabase(cfg *config.Config) *sqlx.DB {
	db, err := sqlx.Connect("postgres",
		"postgres://"+cfg.DBUser+":"+cfg.DBPassword+"@"+cfg.DBHost+":"+cfg.DBPort+"/"+cfg.DBName+"?sslmode=disable")
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	return db
}

func runMigration(db *sqlx.DB, filePath string) error {
	sqlBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading migration file: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing migration: %w", err)
	}

	log.Printf("Migration applied from file: %s", filePath)
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

	log.Printf("Server starting on port %s...", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Error starting the server: %v", err)
	}
}
