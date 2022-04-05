package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ShaghayeghFathi/errandboi/internal/http/request"
	"github.com/ShaghayeghFathi/errandboi/internal/http/response"
	"github.com/ShaghayeghFathi/errandboi/internal/model"
	"github.com/ShaghayeghFathi/errandboi/internal/store/mongo"
	redisp "github.com/ShaghayeghFathi/errandboi/internal/store/redis"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

type Handler struct {
	Redis  *redisp.RedisDB
	Mongo  *mongo.DB
	Logger *zap.Logger
}

func (h Handler) registerEvents(ctx *fiber.Ctx) error {
	action := new(request.Action)

	if err := ctx.BodyParser(action); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	actionID := primitive.NewObjectID()
	actionIDValue := actionID.Hex()
	temp := action.Type[0]

	const t = 2
	if len(action.Type) == t {
		temp = action.Type[0] + "_" + action.Type[1]
	}

	for i := 0; i < len(action.Events); i++ {
		releaseTime := calculateReleaseTime(action.Events[i].Delay)
		id := actionIDValue + "_" + strconv.Itoa(i)

		h.cacheEvent(ctx.Context(), id, releaseTime, action.Events[i], temp)

		eventModel := &model.Event{
			ID: id, ActionID: actionIDValue,
			Description: action.Events[i].Description,
			Delay:       action.Events[i].Delay,
			Topic:       action.Events[i].Topic,
			Payload:     action.Events[i].Payload,
			Status:      "pending",
		}

		_, err := h.Mongo.StoreEvent(ctx.Context(), eventModel)
		if err != nil {
			h.Logger.Error(
				"failed to store event",
				zap.Error(err),
			)
		}
	}

	_, err := h.Mongo.StoreAction(ctx.Context(), actionID, action.Type, len(action.Events))
	if err != nil {
		h.Logger.Error(
			"could not store action",
			zap.Error(err),
		)
	}

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"id": actionIDValue,
	})
}

func (h *Handler) getEvents(ctx *fiber.Ctx) error {
	eventID := ctx.Params("eventId")
	objectID, _ := primitive.ObjectIDFromHex(eventID)
	action, _ := h.Mongo.GetAction(ctx.Context(), objectID)
	events, _ := h.Mongo.GetEvents(ctx.Context(), eventID)

	return ctx.Status(http.StatusOK).JSON(response.GetEventsResponse{Type: action.Type, Events: events})
}

func (h *Handler) getEventStatus(ctx *fiber.Ctx) error {
	eventID := ctx.Params("eventId")

	events, _ := h.Mongo.GetEventStatus(ctx.Context(), eventID)

	s := "done"

	for i := 0; i < len(events); i++ {
		if events[i].Status == "pending" {
			s = "pending"

			break
		}
	}

	return ctx.Status(http.StatusOK).JSON(response.GetEventsStatusResponse{Status: s, Events: events})
}

func (h *Handler) cacheEvent(ctx context.Context, id string, releaseTime float64, event request.Event, temp string) {
	_, err := h.Redis.ZSet(ctx, "events", releaseTime, id)
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

	p, _ := event.Payload.(string)
	err = h.Redis.Set(ctx, "payload"+"_"+id, p)

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
	b := 64
	delayInt, _ := strconv.ParseFloat(delay, b)
	releaseTime := 0.0

	t := 60.0

	switch unit {
	case "s":
		releaseTime = float64(time.Now().Unix()) + delayInt
	case "m":
		releaseTime = float64(time.Now().Unix()) + delayInt*t
	case "h":
		releaseTime = float64(time.Now().Unix()) + delayInt*t*t
	}

	return releaseTime
}

func (h Handler) Register(g fiber.Router) {
	g.Post("", h.registerEvents)
	g.Get("/:eventId", h.getEvents)
	g.Get("/:eventId/status", h.getEventStatus)
}
