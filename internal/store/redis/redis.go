package redis

import (
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"

	"errandboi/internal/db/rdb"
)
type RedisDB struct {
	db *rdb.Redis
}

func (r *RedisDB) ZSet(ctx context.Context, setName string, publishTime float64, id string) (int64 ,error) {
	return r.db.Client.ZAddNX(ctx, setName, &redis.Z{ Score:  publishTime, Member: id}).Result()
}

func (r *RedisDB) ZGetRange(ctx context.Context, setName string, start int, end int) ([]redis.Z, error) {
	return r.db.Client.ZRangeArgsWithScores(ctx, redis.ZRangeArgs{Key: setName,Start: start,Stop:end,ByScore: true,Rev:true}).Result()
}

func (r *RedisDB) ZDel(ctx context.Context, setName string)([]redis.Z , error){
	return r.db.Client.ZPopMin(ctx, setName).Result()
}