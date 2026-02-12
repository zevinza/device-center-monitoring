package cache

import (
	"api/lib"
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// REDIS null if not initialized
var Redis *redis.Client

// InitRedis initialize redis connection
func Connect() {
	if nil == Redis {
		redisHost := viper.GetString("REDIS_HOST")
		redisPort := viper.GetString("REDIS_PORT")
		if redisHost != "" {
			Redis = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
				Password: viper.GetString("REDIS_PASS"),
				DB:       viper.GetInt("REDIS_INDEX"),
			})
		}
	}
}

// SetCachingRedis func
func SetCachingRedis(rdb *redis.Client, datas map[string]map[string]interface{}) {
	repo := NewRedisRepository(rdb)
	re := regexp.MustCompile(`[0-9]+$`)
	for k, v := range datas {
		key := re.ReplaceAll([]byte(k), []byte(``))
		err := repo.Set(string(key), lib.JsonString(v["values"]), 0)
		if err != nil {
			fmt.Printf("unable to SET data. error: %v", err)
		}
	}
}

// RedisRepository represent the repositories
type RedisRepository interface {
	Set(key string, value interface{}, exp time.Duration) error
	Get(key string) (string, error)
	SetStr(key string, value string, exp time.Duration) error
	GetDel(key string) (string, error)
}

// repository represent the repository model
type redisRepository struct {
	Client *redis.Client
}

// NewRedisRepository will create an object that represent the Repository interface
func NewRedisRepository(Client *redis.Client) RedisRepository {
	return &redisRepository{Client}
}

// Set attaches the redis repository and set the data
func (r *redisRepository) Set(key string, value interface{}, exp time.Duration) error {
	return r.Client.Set(context.Background(), key, value, exp).Err()
}

// Get attaches the redis repository and get the data
func (r *redisRepository) Get(key string) (string, error) {
	get := r.Client.Get(context.Background(), key)
	return get.Result()
}

// Set attaches the redis repository and set the data string
func (r *redisRepository) SetStr(key string, value string, exp time.Duration) error {
	return r.Client.Set(context.Background(), key, value, exp).Err()
}

// Delete the data by key
func (r *redisRepository) GetDel(key string) (string, error) {
	get := r.Client.GetDel(context.Background(), key)
	return get.Result()
}
