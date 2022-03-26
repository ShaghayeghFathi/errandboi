package handler

import (
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"errandboi/internal/http/request"
	"errandboi/internal/http/response"

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
	actionIdValue:= actionId.Hex()
	fmt.Println("actionIdValue: ", actionIdValue)
	for i := 0; i < len(action.Events); i++ {
		releaseTime := calculateReleaseTime(action.Events[i].Delay)
		id := actionIdValue + "_" + strconv.Itoa(i)
		h.Redis.ZSet(ctx.Context() , "events" , releaseTime, id ) // TODO: add set name to config

		h.Mongo.StoreEvent(ctx.Context(), id , actionIdValue, action.Events[i].Description, action.Events[i].Delay , 
		action.Events[i].Topic, action.Events[i].Payload)
	}
	h.Mongo.StoreAction(ctx.Context(), actionId, action.Type,len(action.Events))
	
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"id" : actionIdValue,
	  })
}

func(h *Handler) getEvents(ctx *fiber.Ctx)error{
	eventId := ctx.Params("eventId")
	objectId,_ := primitive.ObjectIDFromHex(eventId)
	action,_ := h.Mongo.GetAction(ctx.Context(), objectId)
	events,_ := h.Mongo.GetEvents(ctx.Context(), eventId)	
	return ctx.Status(http.StatusOK).JSON(response.GetEventsResponse{Type: action.Type, Events:events})
}

func(h *Handler) getEventStatus(ctx *fiber.Ctx)error{
	eventId := ctx.Params("eventId")
	objectId,_ := primitive.ObjectIDFromHex(eventId)
	action,_ := h.Mongo.GetAction(ctx.Context(), objectId)
	events,_ := h.Mongo.GetEventStatus(ctx.Context(), eventId)	
	return ctx.Status(http.StatusOK).JSON(response.GetEventsStatusResponse{Status: action.Status, Events:events})
}

func calculateReleaseTime(delaySt string) float64{
	unit := string(delaySt[len(delaySt)-1])
	delay := delaySt[:len(delaySt)-1]
	delayInt,_ := strconv.ParseFloat(delay, 64)
	releaseTime :=0.0
	switch unit {
	case "s":{
		releaseTime =float64(time.Now().Unix())+delayInt
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
	app.Get("/events/:eventId", h.getEvents)
	app.Get("/events/:eventId/status", h.getEventStatus)	
}
