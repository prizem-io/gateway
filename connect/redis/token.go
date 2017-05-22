package redis

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	"gopkg.in/vmihailenco/msgpack.v2"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/utils"
)

type (
	redisClient interface {
		SetNX(key string, value interface{}, expiration time.Duration) *redis.BoolCmd
		Get(key string) *redis.StringCmd
		Expire(key string, expiration time.Duration) *redis.BoolCmd
	}

	RedisTokener struct {
		redis redisClient
	}
)

func NewTokener(redis redisClient) *RedisTokener {
	return &RedisTokener{
		redis: redis,
	}
}

func (d *RedisTokener) Create(token *config.Token) (*config.Token, error) {
	set := false
	for tries := 3; !set && tries >= 0; tries-- {
		token.ID = utils.RandomString(64)
		bytes, err := msgpack.Marshal(token)
		if err != nil {
			return nil, err
		}

		var timeout time.Duration
		if token.Expiry > 0 {
			timeout = time.Duration(token.Expiry) * time.Second
		}

		set, err = d.redis.SetNX(token.ID, bytes, timeout).Result()
		if err != nil {
			return nil, err
		}
	}
	if !set {
		return nil, errors.New("Could not create token")
	}

	return token, nil
}

func (d *RedisTokener) Get(id string) (*config.Token, error) {
	data, err := d.redis.Get(id).Bytes()
	if err != nil {
		return nil, err
	}

	token := config.Token{}
	err = msgpack.Unmarshal(data, &token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (d *RedisTokener) Touch(token *config.Token) error {
	var timeout time.Duration
	if token.Expiry > 0 {
		timeout = time.Duration(token.Expiry) * time.Second
	}

	return d.redis.Expire(token.ID, timeout).Err()
}
