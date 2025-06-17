package postgres

import (
	"database/sql"
	"fmt"

	configloader "github.com/cafpleon/filingo-config"
	_ "github.com/jackc/pgx/v5/stdlib" // El driver de pgx para database/sql
)

func Connect(cfg configloader.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("falló al abrir la conexión a postgres: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("falló el ping a postgres: %w", err)
	}

	// Aquí podrías configurar el pool de conexiones:
	// db.SetMaxOpenConns(...)

	return db, nil
}
