package handler

import (
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"errandboi/internal/http/request"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type Handler struct{
	Redis *redisPK.RedisDB
}

func(h Handler) registerEvents(ctx *fiber.Ctx)error{
	action := new(request.Action)

	if err := ctx.BodyParser(action); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	actionId := uuid.NewString()

	for i := 0; i < len(action.Events); i++ {
		releaseTime := calculateReleaseTime(action.Events[i].Delay)
		id := actionId + "_" + strconv.Itoa(i)
		h.Redis.ZSet(ctx.Context() , "events" , releaseTime, id ) // TODO: add set name to config
	}
	
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"id" : actionId,
	  })
}

func calculateReleaseTime(delaySt string) float64{
	unit := string(delaySt[len(delaySt)-1])
	delay := delaySt[:len(delaySt)-1]
	delayInt,_ := strconv.ParseFloat(delay, 64)
	releaseTime :=0.0
	switch unit {
	case "s":{
		releaseTime =float64(time.Now().Unix())+delayInt
		fmt.Println("time unix : " , float64(time.Now().Unix()))
	}
	case "m":{
		releaseTime = float64(time.Now().Unix())+ delayInt*60
	}
	case "h":
		releaseTime = float64(time.Now().Unix())+ delayInt*60*60
	}
	return releaseTime
}

func (h Handler) Register(app *fiber.App) {
	app.Post("/events", h.registerEvents)
	
}
