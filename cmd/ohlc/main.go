package main

import (
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"

	"github.com/dneprix/ohlc/pkg/downloaders"
	"github.com/dneprix/ohlc/pkg/schedulers"
)

func main() {
	// Setup logger
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// DB connection
	dbConnection := os.Getenv("DB_CONNECTION")
	db, err := sqlx.Open("postgres", dbConnection)
	if err != nil {
		logger.Fatal(err)
	}
	db.SetMaxIdleConns(50)
	db.SetConnMaxLifetime(time.Second * 5)
	defer db.Close()

	// DB migrate
	dbMigrationsPath := os.Getenv("DB_MIGRATIONS_PATH")
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(dbMigrationsPath, "postgres", driver)
	if err != nil {
		logger.Fatal(err)
	}
	m.Up()

	// Initialise scheduler
	scheduler := schedulers.NewSheduler(logger)

	// Add downloaders to scheduler
	scheduler.Add(downloaders.NewCryptowatDownloader(db, logger))
	scheduler.Add(downloaders.NewKrakenDownloader(db, logger))
	scheduler.Add(downloaders.NewCryptocompareDownloader(db, logger))

	// Run scheduler
	scheduler.Run()
}
