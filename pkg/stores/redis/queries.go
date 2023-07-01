package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

// Options represents the cache store available options
type Options struct {
	// Cost corresponds to the memory capacity used by the item when setting a value
	// Actually it seems to be used by Ristretto library only
	Cost int64

	// Expiration allows to specify an expiration time when setting a value
	Expiration time.Duration

	// Tags allows to specify associated tags to the current value
	Tags []string
}

// CostValue returns the allocated memory capacity
func (o Options) CostValue() int64 {
	return o.Cost
}

// ExpirationValue returns the expiration option value
func (o Options) ExpirationValue() time.Duration {
	return o.Expiration
}

// TagsValue returns the tags option value
func (o Options) TagsValue() []string {
	return o.Tags
}

type Cache struct {
	client *redis.Client
}

func NewCache() *Cache {
	return &Cache{
		client: RedisClient(),
	}
}

/*
	Set the key, value data to the redis cache.
	Args:
		ctx context.Context
		key string
		value interface{}
		options Options struct type
	options.ExpirationValue() gives the ttl for the key stored in the redis cache.
*/
func (c *Cache) Set(ctx context.Context, key string, value interface{}, options *Options) {
	if options == nil {
		options = &Options{}
	}
	err := c.client.Set(ctx, key, value, options.ExpirationValue()).Err()
	if err != nil {
		log.Println("Redis Set error", err)
	}
}

/*
	Get the data from the redis cache
	Args:
		ctx context.Context
		key string
	Returns the data present for the key in redis cache
*/
func (c *Cache) Get(ctx context.Context, key string) interface{} {
	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		return nil
	}
	return val
}

/*
	Delete the data from redis cache based on key
	Args:
		ctx context.Context
		key string
*/
func (c *Cache) Delete(ctx context.Context, key string) {
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		log.Println("Delete Cache error", err)
	}
}

/*
	Delete the data from redis cache based on matching of the key
	Args:
		ctx context.Context
		key string
*/
func (c *Cache) DeleteMatching(ctx context.Context, key string) {
	iter := c.client.Scan(ctx, 0, key, 0).Iterator()
	for iter.Next(ctx) {
		c.client.Del(ctx, iter.Val())
	}
}

func (c *Cache) GetAllMatchingKeys(ctx context.Context, key string) *redis.ScanIterator {
	iter := c.client.Scan(ctx, 0, key, 0).Iterator()
	return iter
}

/*
InvalidateCache Generic method to invalidate the cache.
Args:
	ctx contex.Context
	sender string For eg. canvas, user
	uuid []string
	bulk bool
*/
func InvalidateCache(ctx context.Context, sender string, uuid []string, bulk bool) {
	key := GenerateCacheKey(sender, uuid)
	cache := NewCache()
	if bulk {
		key += ":*"
		cache.DeleteMatching(ctx, key)
	} else {
		cache.Delete(ctx, key)
	}
}

/*
	HSet the key, value data to the redis cache.
	Args:
		ctx context.Context
		hash string
		key string
		value interface{}
		options Options struct type
	options.ExpirationValue() gives the ttl for the key stored in the redis cache.
*/
func (c *Cache) HSet(ctx context.Context, hash string, key string, value interface{}) {
	fmt.Println(hash)
	fmt.Println(key)
	err := c.client.HSet(ctx, hash, key, value).Err()
	if err != nil {
		log.Println("Redis HSet error", err)
	}
}

/*
	Get the data from the redis cache
	Args:
		ctx context.Context
		hash string
		key string
	Returns the data present for the key in redis cache
*/
func (c *Cache) HGet(ctx context.Context, hash string, key string) string {
	res := c.client.HGet(ctx, hash, key)
	return res.Val()
}

/*
	HDelete the data from the redis cache
	Args:
		ctx context.Context
		hash string
		key string
	Returns the data present for the key in redis cache
*/
func (c *Cache) HDelete(ctx context.Context, hash string, key string) int64 {
	res := c.client.HDel(ctx, hash, key)
	fmt.Println("redis response delete", res.Val())
	return res.Val()
}

func (c *Cache) Hscan(ctx context.Context, hash string, keyPattern string) []string {
	res := c.client.HScan(ctx, hash, 0, keyPattern, 100000)
	val, _ := res.Val()
	return val
}

func (c *Cache) HDeleteMatching(ctx context.Context, hash string, keyPattern string) {
	values := c.Hscan(ctx, hash, keyPattern)
	for i, key := range values {
		if i%2 == 0 {
			c.HDelete(ctx, hash, key)
		}
	}
}

func (c *Cache) HGetAll(ctx context.Context, keyPattern string) map[string]string {
	res := c.client.HGetAll(ctx, keyPattern)
	return res.Val()
}
