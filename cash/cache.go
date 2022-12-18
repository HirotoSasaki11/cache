package cash

import (
	"context"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	lru "github.com/hashicorp/golang-lru"
)

type Cacher interface {
	Init(csh *Cash)

	LocalCache() bool

	Load(key string) ([]byte, error)

	Store(key string, value []byte) error

	Delete(key string) error
}

type nilCacher struct{}

func (c *nilCacher) Init(csh *Cash) {
}

func (c *nilCacher) LocalCache() bool {
	return true
}

func (c *nilCacher) Load(key string) ([]byte, error) {
	return nil, ErrNoCache
}

func (c *nilCacher) Store(key string, value []byte) error {
	return nil
}

func (c *nilCacher) Delete(key string) error {
	return nil
}

type LRUCache struct {
	Size int

	csh *Cash

	c *lru.Cache
}

func (c *LRUCache) Init(csh *Cash) {
	var err error
	c.c, err = lru.New(c.Size)
	if err != nil {
		panic(err)
	}
	c.csh = csh
}

func (c *LRUCache) LocalCache() bool {
	return true
}

func (c *LRUCache) Load(key string) ([]byte, error) {
	v, ok := c.c.Get(key)
	if !ok {
		return nil, ErrNoCache
	}
	return v.([]byte), nil
}

func (c *LRUCache) Store(key string, value []byte) error {
	c.c.Add(key, value)
	return nil
}

func (c *LRUCache) Delete(key string) error {
	c.c.Remove(key)
	return nil
}

type MapCache struct {
	csh *Cash

	c sync.Map
}

func (c *MapCache) Init(csh *Cash) {
	c.csh = csh
}

func (c *MapCache) LocalCache() bool {
	return true
}

func (c *MapCache) Load(key string) ([]byte, error) {
	v, ok := c.c.Load(key)
	if !ok {
		return nil, ErrNoCache
	}
	return v.([]byte), nil
}

func (c *MapCache) Store(key string, value []byte) error {
	c.c.Store(key, value)
	return nil
}

func (c *MapCache) Delete(key string) error {
	c.c.Delete(key)
	return nil
}

const defaultRedisTTL time.Duration = 7 * 24 * time.Hour

type RedisCache struct {
	Pool *redis.Pool

	ConnectTimeout time.Duration

	TTL time.Duration

	KeyPrefix string

	csh *Cash
}

func (c *RedisCache) Init(csh *Cash) {
	c.csh = csh
	if c.KeyPrefix == "" {
		c.KeyPrefix = "cash:"
	}
}

func (c *RedisCache) LocalCache() bool {
	return false
}

func (c *RedisCache) getConn() (redis.Conn, error) {
	var zeroDuration time.Duration
	ctx := context.Background()
	if c.ConnectTimeout != zeroDuration {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.ConnectTimeout)
		defer cancel()
	}
	return c.Pool.GetContext(ctx)
}

func (c *RedisCache) Load(key string) ([]byte, error) {
	conn, err := c.getConn()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	v, err := redis.Bytes(conn.Do("GET", c.KeyPrefix+key))
	if err == redis.ErrNil || len(v) == 0 {
		return nil, ErrNoCache
	} else if err != nil {
		return nil, err
	} else {
		return v, nil
	}
}

func (c *RedisCache) Store(key string, value []byte) error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	ttl := c.TTL
	if ttl == 0 {
		ttl = defaultRedisTTL
	}

	_, err = conn.Do("SET", c.KeyPrefix+key, value, "EX", ttl.Seconds())
	return err
}

func (c *RedisCache) Delete(key string) error {
	conn, err := c.getConn()
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Do("DEL", c.KeyPrefix+key)
	return err
}
