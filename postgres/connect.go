package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	// Ahora solo importamos pgxpool directamente
	"github.com/jackc/pgx/v5/pgxpool"

	configloader "github.com/cafpleon/filingo-util-config" // Tu librería de config
)

// Connect crea y devuelve un nuevo y performante pool de conexiones de pgx.
// La firma ahora devuelve el tipo específico *pgxpool.Pool.
func Connect(ctx context.Context, cfg configloader.DBConfig) (*pgxpool.Pool, error) {
	// 1. Construir la URL de conexión (DSN)
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)
	slog.Info("Intentando conectar a la base de datos con pgxpool", "host", cfg.Host, "db", cfg.Name)

	// 2. Parsear la configuración del pool a partir del DSN.
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("no se pudo parsear la configuración del pool de pgx: %w", err)
	}

	// 3. Aplicar la configuración específica del pool.
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifeTime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Se puede añadir un timeout al contexto para la conexión inicial.
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 4. Crear el pool de conexiones usando directamente pgxpool.
	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el pool de conexiones: %w", err)
	}

	// 5. Verificar que la conexión está viva.
	pingCtx, cancelPing := context.WithTimeout(ctx, 3*time.Second)
	defer cancelPing()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	slog.Info("Conexión a la base de datos (pgxpool) establecida exitosamente.")
	return pool, nil
}
