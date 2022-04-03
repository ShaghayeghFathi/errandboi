package handler

import (
	"context"
	"errandboi/internal/model"
	"errandboi/internal/store/mongo"
	redisPK "errandboi/internal/store/redis"
	"log"
	"net/http"
	"strconv"
	"time"

	"errandboi/internal/http/request"
	"errandboi/internal/http/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Handler struct {
	Redis  *redisPK.RedisDB
	Mongo  *mongo.MongoDB
	Logger *zap.Logger
}

func (h Handler) registerEvents(ctx *fiber.Ctx) error {
	action := new(request.Action)

	if err := ctx.BodyParser(action); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	actionId := primitive.NewObjectID()
	actionIdValue := actionId.Hex()
	temp := action.Type[0]

	if len(action.Type) == 2 {
		temp = action.Type[0] + "_" + action.Type[1]
	}

	for i := 0; i < len(action.Events); i++ {
		releaseTime := calculateReleaseTime(action.Events[i].Delay)
		id := actionIdValue + "_" + strconv.Itoa(i)

		h.cacheEvent(ctx.Context(), id, releaseTime, action.Events[i], temp)

		eventModel := &model.Event{ID: id, ActionId: actionIdValue,
			Description: action.Events[i].Description,
			Delay:       action.Events[i].Delay,
			Topic:       action.Events[i].Topic,
			Payload:     action.Events[i].Payload,
			Status:      "pending"}

		_, err := h.Mongo.StoreEvent(ctx.Context(), eventModel)
		if err != nil {
			h.Logger.Error(
				"failed to store event",
				zap.Error(err),
			)
		}
	}

	_, err := h.Mongo.StoreAction(ctx.Context(), actionId, action.Type, len(action.Events))
	if err != nil {
		h.Logger.Error(
			"could not store action",
			zap.Error(err),
		)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"id": actionIdValue,
	})
}

func (h *Handler) getEvents(ctx *fiber.Ctx) error {
	eventId := ctx.Params("eventId")
	objectId, _ := primitive.ObjectIDFromHex(eventId)
	action, _ := h.Mongo.GetAction(ctx.Context(), objectId)
	events, _ := h.Mongo.GetEvents(ctx.Context(), eventId)
	return ctx.Status(http.StatusOK).JSON(response.GetEventsResponse{Type: action.Type, Events: events})
}

func (h *Handler) getEventStatus(ctx *fiber.Ctx) error {

	eventId := ctx.Params("eventId")

	events, _ := h.Mongo.GetEventStatus(ctx.Context(), eventId)

	var s string = "done"
	for i := 0; i < len(events); i++ {
		if events[i].Status == "pending" {
			s = "pending"
			break
		}
	}

	return ctx.Status(http.StatusOK).JSON(response.GetEventsStatusResponse{Status: s, Events: events})
}

func (h *Handler) cacheEvent(ctx context.Context, id string, releaseTime float64, event request.Event, temp string) {
	_, err := h.Redis.ZSet(ctx, "events", releaseTime, id) // TODO: add set name to config
	if err != nil {
		log.Fatal(err)
	}
	err = h.Redis.Set(ctx, "desc"+"_"+id, event.Description)
	if err != nil {
		log.Fatal(err)
	}
	err = h.Redis.Set(ctx, "topic"+"_"+id, event.Topic)
	if err != nil {
		log.Fatal(err)
	}
	err = h.Redis.Set(ctx, "payload"+"_"+id, event.Payload.(string))
	if err != nil {
		log.Fatal(err)
	}
	err = h.Redis.Set(ctx, "type"+"_"+id, temp)
	if err != nil {
		log.Fatal(err)
	}
}

func calculateReleaseTime(delaySt string) float64 {

	unit := string(delaySt[len(delaySt)-1])
	delay := delaySt[:len(delaySt)-1]
	delayInt, _ := strconv.ParseFloat(delay, 64)
	releaseTime := 0.0

	switch unit {
	case "s":
		releaseTime = float64(time.Now().Unix()) + delayInt
	case "m":
		releaseTime = float64(time.Now().Unix()) + delayInt*60
	case "h":
		releaseTime = float64(time.Now().Unix()) + delayInt*60*60
	}
	return releaseTime
}

func (h Handler) Register(app *fiber.App) {
	app.Post("/events", h.registerEvents)
	app.Get("/events/:eventId", h.getEvents)
	app.Get("/events/:eventId/status", h.getEventStatus)
}
