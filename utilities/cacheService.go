package utilities

import "github.com/go-redis/redis"

type CacheService struct {
	redisClient *redis.Client
}

func NewCacheService(redisClient *redis.Client) *CacheService {
	return &CacheService{redisClient}
}

func (cs *CacheService) Get(key string) (bool, string, error) {
	result, err := cs.redisClient.Get(key).Result()
	if err == redis.Nil {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}
	return true, result, nil
}

func (cs *CacheService) Set(key string, value interface{}) {
	cs.redisClient.Set(key, value, 0)
}

func (cs *CacheService) Remove(key string) {
	cs.redisClient.Del(key)
}
