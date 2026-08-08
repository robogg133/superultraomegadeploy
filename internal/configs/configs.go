package configs

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/jackc/pgx/v5/pgxpool"
)

const deleteOperation = "DELETE"

const (
	KeyRegistryURL      string = "registryUrl"
	KeyRegistryUser     string = "registryUser"
	KeyRegistryPassword string = "registryPassword"
)

type ConfigValue struct {
	UpdatedAt time.Time
	Value     json.RawMessage
}

type configChange struct {
	UpdatedAt time.Time       `json:"updated_at"`
	Value     json.RawMessage `json:"value"`
	Operation string          `json:"operation"`
	Key       string          `json:"key"`
}

type cfgs struct {
	pool   *pgxpool.Pool
	mu     sync.RWMutex
	values map[string]ConfigValue
}

var Configs *cfgs

func init() {
	Configs = new(cfgs)
	Configs.values = make(map[string]ConfigValue)
}

func Init(ctx context.Context, pool *pgxpool.Pool) error {
	Configs.pool = pool

	rows, err := pool.Query(ctx, `SELECT key, value, updated_at FROM configs.configs`)
	if err != nil {
		return err
	}
	defer rows.Close()

	Configs.mu.Lock()
	defer Configs.mu.Unlock()

	for rows.Next() {
		var key string
		var v ConfigValue
		if err := rows.Scan(&key, &v.Value, &v.UpdatedAt); err != nil {
			return err
		}
		Configs.values[key] = ConfigValue{
			Value:     v.Value,
			UpdatedAt: v.UpdatedAt,
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return Configs.listenChange(ctx)
}

// listenChange creates a goroutine to update Configs map
func (c *cfgs) listenChange(ctx context.Context) error {
	conn, err := c.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer conn.Release()

	if _, err := conn.Exec(ctx, "LISTEN config_changes"); err != nil {
		return err
	}

	go func() {
		for {
			l := log.Ctx(ctx)

			noti, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				l.Err(err).Send()
				time.Sleep(time.Second)
				continue
			}

			var change configChange
			if err := json.Unmarshal([]byte(noti.Payload), &change); err != nil {
				l.Err(err).Send()
				continue
			}

			Configs.mu.Lock()
			if change.Operation == deleteOperation {
				delete(Configs.values, change.Key)
			} else {
				Configs.values[change.Key] = ConfigValue{
					UpdatedAt: change.UpdatedAt,
					Value:     change.Value,
				}
			}
			Configs.mu.Unlock()
			log.Info().Str("key", change.Key).Str(string(change.Value), "value").Msg("config update")
		}
	}()
	return nil
}
