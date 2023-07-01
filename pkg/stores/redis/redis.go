package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
)

// Declaring the ctx, rdb values. rdb value is used globaly to access redis db methods.
var rdb *redis.Client

// Initializing the redis cache.
// Redis config keys we are taking from config.yaml file
// InitRedis method is triggered in the main method.
func InitRedis() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", configs.GetRedisConfig().Host, configs.GetRedisConfig().Port),
		Password: configs.GetRedisConfig().Password,
		DB:       0,
	})
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Println("Error on connecting to Redis: ", err)
	}
	log.Println("Connected to Redis: ", pong)
	return rdb
}

// RedisClient return the rdb(redis database) client instance.
func RedisClient() *redis.Client {
	return rdb
}
func GetBgContext() context.Context {
	return context.Background()

}
func SetHash(key string, value map[string]interface{}) {
	// Create arguments: key field value [field value]...
	var args = []interface{}{key}
	for k, v := range value {
		args = append(args, k, v)
	}
	rdb.Do(GetBgContext(), args...)

	//_, err := conn.Do("HMSET", args...)
	//if err != nil {
	//	return fmt.Errorf("error setting key %s to %v: %v", key, value, err)
	//}

}

//
//func SetKey(c *redis.Client, key string, value interface{}) error {
//	p, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	return c.Set(key, p)
//}
