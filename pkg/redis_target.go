package pkg

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// RedisTarget is a struct that holds the data for the Redis target
type RedisTarget struct {
	Host     string
	Port     string
	Password string
}

// NewRedisTarget creates a new RedisTarget struct
func NewRedisTarget() *RedisTarget {
	return &RedisTarget{}
}

// WriteToRedis writes the data from the CSVSource struct to the Redis target
func (r *RedisTarget) WriteToRedis(source *CSVSource, prefix string) error {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Host + ":" + r.Port,
		Password: r.Password,
		DB:       0,
	})
	ping_cmd := client.Ping(ctx)
	if ping_cmd.Err() != nil {
		return ping_cmd.Err()
	}

	for ui, row := range source.Rows {
		for i, value := range row {
			keyname := fmt.Sprintf("%s%s-%v", prefix, source.Header[i], ui)
			err := client.Set(ctx, keyname, value, 0).Err()
			if err != nil {
				return err
			}
			fmt.Printf("SET %s %s\n", keyname, value)
		}
	}

	return nil
}
