package dbconnector

import (
	"database/sql"
	"fmt"

	configloader "github.com/cafpleon/filingo-util-config"   // Tu librería de config
	"github.com/cafpleon/filingo-util-db-connector/mysql"    // Importa tu sub-paquete mysql
	"github.com/cafpleon/filingo-util-db-connector/postgres" // Importa tu sub-paquete postgres
)

// Connect es la función principal y pública. Actúa como una fábrica que
// elige el conector correcto basándose en la configuración.
func Connect(cfg configloader.DBConfig) (*sql.DB, error) {
	switch cfg.Driver {
	case "postgres":
		return postgres.Connect(cfg)
	case "mysql":
		return mysql.Connect(cfg)
	default:
		return nil, fmt.Errorf("driver de base de datos no soportado: %s", cfg.Driver)
	}
}
