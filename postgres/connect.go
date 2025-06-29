package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5" // Lo usamos solo por el error pgx.ErrNoRows
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	configloader "github.com/cafpleon/filingo-util-config" // Tu librería de config
)

// IDBPool define el contrato mínimo que cualquier repositorio necesitará
// para interactuar con el pool de conexiones de pgx.
type IDBPool interface {
	// Método de pgx para consultas que devuelven múltiples filas
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)

	// Método de pgx para consultas que devuelven una sola fila
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row

	// Método de pgx para operaciones de escritura (INSERT, UPDATE, DELETE)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)

	//  Begin inicia una transacción. El objeto Tx también cumple la interfaz IDBPool
	// (sin los métodos Begin y Close), por lo que podemos reutilizarlo.
	Begin(ctx context.Context) (pgx.Tx, error)

	// Método para cerrar el pool
	Close()
}

// Connect crea y devuelve un nuevo y performante pool de conexiones de pgx.
// La firma ahora devuelve el tipo específico *pgxpool.Pool.
func ConnectToPostgres(ctx context.Context, cfg configloader.DBConfig) (*pgxpool.Pool, error) {
	//  Construir la URL de conexión (DSN)
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)
	slog.Info("Intentando conectar a la base de datos con pgxpool", "host", cfg.Host, "db", cfg.Name)

	// Parsear la configuración del pool a partir del DSN.
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("no se pudo parsear la configuración del pool de pgx: %w", err)
	}

	// Aplicar la configuración específica del pool.
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MinConns = cfg.MinConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifeTime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.HealthCheckPeriod

	// Se puede añadir un timeout al contexto para la conexión inicial.
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Crear el pool de conexiones usando directamente pgxpool.
	pool, err := pgxpool.NewWithConfig(connectCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el pool de conexiones: %w", err)
	}

	// Verificar que la conexión está viva.
	pingCtx, cancelPing := context.WithTimeout(ctx, 3*time.Second)
	defer cancelPing()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo hacer ping a la base de datos: %w", err)
	}

	slog.Info("Conexión a la base de datos (pgxpool) establecida exitosamente.")
	return pool, nil
}
