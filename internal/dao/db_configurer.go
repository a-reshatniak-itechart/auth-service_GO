package dao

import (
	"auth-service/internal/config"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	migratePostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func Configure(config *config.Config) *gorm.DB {
	var dbHost string
	if os.Getenv("POSTGRES_HOST") == "" {
		dbHost = config.DBHost
	} else {
		dbHost = os.Getenv("POSTGRES_HOST")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", dbHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Println("? Connected Successfully to the Database")
	//to fix 'Function uuid_generate_v4() does not exist'
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Migrate the schema
	//db.AutoMigrate(&internal.User{})
	migrateDb(config)

	return db
}

func migrateDb(config *config.Config) {
	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBPort, config.DBName)
	db, _ := sql.Open("postgres", databaseUrl)
	driver, _ := migratePostgres.WithInstance(db, &migratePostgres.Config{})
	m, _ := migrate.NewWithDatabaseInstance(config.MigrationFolder, config.DBName, driver)
	m.Up()
}
