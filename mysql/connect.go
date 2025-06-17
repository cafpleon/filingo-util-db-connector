package mysql

import (
	"database/sql"
	"fmt"

	configloader "github.com/cafpleon/filingo-util-config"
	_ "github.com/go-sql-driver/mysql" // El driver de MySQL
)

func Connect(cfg configloader.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("falló al abrir la conexión a mysql: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("falló el ping a mysql: %w", err)
	}

	return db, nil
}
