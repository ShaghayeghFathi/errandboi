package redisPK

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"errandboi/internal/db/rdb"
)
type RedisDB struct {
	db *rdb.Redis
}

func NewRedis(db *rdb.Redis) *RedisDB{
	return &RedisDB{db:db}
}

func (r *RedisDB) ZSet(ctx context.Context, setName string, publishTime float64, id string) (int64 ,error) {
	return r.db.Client.ZAddNX(ctx, setName, &redis.Z{ Score:  publishTime, Member: id}).Result()
}

func (r *RedisDB) ZGetRange(ctx context.Context, setName string, start float64, end float64) ([]redis.Z, error) {
	return r.db.Client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{Key: setName,Start: start,Stop:end,ByScore: true,Rev:true}).Result()
}

func (r *RedisDB) ZDel(ctx context.Context, setName string)([]redis.Z , error){
	return r.db.Client.ZPopMin(ctx, setName).Result()
}

func (r *RedisDB) Set(ctx context.Context, key string, value string) error {
	return r.db.Client.Set(ctx, key, value, 0).Err()
}

func (r *RedisDB) Get(ctx context.Context, key string) (string, error) {
	return r.db.Client.Get(ctx, key).Result()
}

func (r *RedisDB) Del(ctx context.Context, key string) error {
	return r.db.Client.Del(ctx, key).Err()
}