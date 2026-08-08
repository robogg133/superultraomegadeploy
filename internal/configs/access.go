package configs

import (
	"context"
	"encoding/json"
	"time"
)

func (c *cfgs) Get(key string) (ConfigValue, bool) {
	c.mu.RLock()
	v, ok := c.values[key]
	c.mu.RUnlock()
	if !ok {
		if df, ok := defaults[key]; ok {
			return ConfigValue{
				UpdatedAt: time.Now(),
				Value:     df,
			}, true
		}
		return ConfigValue{}, false
	}
	return v, true
}

func (c *cfgs) Set(ctx context.Context, key string, val json.RawMessage) error {
	_, err := c.pool.Exec(ctx, `CALL set_config($1, $2)`, key, val)
	return err
}
