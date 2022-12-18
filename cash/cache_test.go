package cash

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/gomodule/redigo/redis"

	"github.com/stretchr/testify/assert"
)

func redisCacherOrErr() (*RedisCache, error) {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			return redis.DialURL("redis://localhost")
		},
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()
	conn, err := pool.GetContext(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := conn.Do("PING"); err != nil {
		return nil, err
	}
	defer conn.Close()
	return &RedisCache{
		Pool:           pool,
		ConnectTimeout: 100 * time.Millisecond,
	}, nil
}

func testCacher(t *testing.T, cacher Cacher) {
	var (
		b   []byte
		err error
	)

	err = cacher.Delete("testing")
	assert.NoError(t, err, "Delete")

	b, err = cacher.Load("testing")
	assert.Empty(t, b)
	assert.EqualError(t, err, ErrNoCache.Error())

	err = cacher.Store("testing", []byte("hi"))
	assert.NoError(t, err, "Store")

	b, err = cacher.Load("testing")
	assert.Equal(t, []byte("hi"), b, "Load")
	assert.NoError(t, err)

	err = cacher.Delete("testing")
	assert.NoError(t, err, "Delete")

	b, err = cacher.Load("testing")
	assert.Empty(t, b)
	assert.EqualError(t, err, ErrNoCache.Error())
}

func TestCacher(t *testing.T) {
	t.Run(fmt.Sprintf("%T", new(LRUCache)), func(t *testing.T) {
		csh := &Cash{}
		cacher := &LRUCache{
			Size: 10,
		}
		cacher.Init(csh)
		testCacher(t, cacher)
	})
	t.Run(fmt.Sprintf("%T", new(MapCache)), func(t *testing.T) {
		csh := &Cash{}
		cacher := &MapCache{}
		cacher.Init(csh)
		testCacher(t, cacher)
	})
	t.Run(fmt.Sprintf("%T", new(RedisCache)), func(t *testing.T) {
		csh := &Cash{}
		cacher, err := redisCacherOrErr()
		if err != nil {
			t.Skip(err)
		}
		cacher.Init(csh)
		testCacher(t, cacher)
	})
}
