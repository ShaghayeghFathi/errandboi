package redispk

import (
	"github.com/ShaghayeghFathi/errandboi/internal/db/rdb"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
)

type RedisDB struct {
	db *rdb.Redis
}

func NewRedis(db *rdb.Redis) *RedisDB {
	return &RedisDB{db: db}
}

func (r *RedisDB) ZSet(ctx context.Context, setName string, publishTime float64, id string) (int64, error) {
	// nolint:wrapcheck
	return r.db.Client.ZAddNX(ctx, setName, &redis.Z{Score: publishTime, Member: id}).Result()
}

func (r *RedisDB) ZGetRange(ctx context.Context, setName string, start float64, end float64) ([]redis.Z, error) {
	// nolint:wrapcheck
	return r.db.Client.ZRangeArgsWithScores(ctx,
		redis.ZRangeArgs{
			Key:     setName,
			Start:   start,
			Stop:    end,
			ByScore: true,
			Rev:     true,
		}).Result()
}

func (r *RedisDB) ZDel(ctx context.Context, setName string) ([]redis.Z, error) {
	// nolint:wrapcheck
	return r.db.Client.ZPopMin(ctx, setName).Result()
}

func (r *RedisDB) ZRem(ctx context.Context, setName string, key string) (int64, error) {
	// nolint:wrapcheck
	return r.db.Client.ZRem(ctx, setName, key).Result()
}

func (r *RedisDB) Set(ctx context.Context, key string, value string) error {
	// nolint:wrapcheck
	return r.db.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisDB) SetInt(ctx context.Context, key string, value int) error {
	// nolint:wrapcheck
	return r.db.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	// nolint:wrapcheck
	return r.db.Client.Get(ctx, key).Result()
}

func (r *RedisDB) Decrement(ctx context.Context, key string) error {
	// nolint:wrapcheck
	return r.db.Client.Decr(ctx, key).Err()
}

// not used

func (r *RedisDB) Del(ctx context.Context, key string) error {
	// nolint:wrapcheck
	return r.db.Client.Del(ctx, key).Err()
}
