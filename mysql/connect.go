package mysql

import (
	"database/sql"
	"fmt"

	configloader "github.com/cafpleon/filingo-util-config"
	_ "github.com/go-sql-driver/mysql" // El driver de MySQL
)

func ConnectToMysql(cfg configloader.DBConfig) (*sql.DB, error) {
	//  Construir la URL de conexión (DSN)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

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
