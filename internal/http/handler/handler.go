package handler

import "errandboi/internal/store/redis"

type Handler struct{
	Redis redis.RedisDB
}

func (h Handler) Register() {

}
