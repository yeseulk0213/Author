package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/golang/glog"
	"strconv"
)

type RedisDB struct {
	client *redis.Client
	ctx context.Context
}

func ConnRedis(ctx context.Context) *RedisDB {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		MinIdleConns: 5,
		PoolSize: 10,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil
	}

	return &RedisDB{
		client: rdb, ctx: ctx,
	}
}

func (r *RedisDB) Set(key string, value interface{}) (string, error) {
	_, err := r.client.Set(r.ctx, key, value, 0).Result()
	return key, err
}

func (r *RedisDB) Get(key string, resultType string) (interface{}, error) {
	result, err := r.client.Get(r.ctx, key).Result()
	if err != nil {
		return result, err
	}

	switch resultType {
	case "uint":
		temp, err := strconv.ParseUint(result, 10, 32)
		if err != nil {
			glog.Fatal(err)
			return temp, err
		}
		return uint(temp), nil
	default:
		return result, err
	}
}

func (r *RedisDB) Delete(key string) (string, error) {
	_, err := r.client.Del(r.ctx, key).Result()
	return key, err
}