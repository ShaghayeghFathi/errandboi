package handler

import (
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"errandboi/internal/http/request"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct{
	Redis *redisPK.RedisDB
	Mongo *mongo.MongoDB
}

func(h Handler) registerEvents(ctx *fiber.Ctx)error{
	action := new(request.Action)

	if err := ctx.BodyParser(action); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	actionId := primitive.NewObjectID()

	for i := 0; i < len(action.Events); i++ {
		releaseTime := calculateReleaseTime(action.Events[i].Delay)
		id := actionId.String() + "_" + strconv.Itoa(i)
		h.Redis.ZSet(ctx.Context() , "events" , releaseTime, id ) // TODO: add set name to config

		h.Mongo.StoreEvent(ctx.Context(), id, action.Events[i].Description, action.Events[i].Delay , action.Events[i].Topic, 
		action.Events[i].Payload)
	}
	h.Mongo.StoreAction(ctx.Context(), actionId, action.Type,len(action.Events))

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"id" : actionId.String(),
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
