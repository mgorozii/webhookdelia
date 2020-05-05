package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mgorozii/webhookdelia/pkg/services/bot"
	"github.com/mgorozii/webhookdelia/pkg/services/store"
)

type service struct {
	*gin.Engine
	store store.Service
	log   *zap.Logger
	bot   bot.Service
}

func (s *service) init(opts Opts) {
	if opts.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}

	s.Engine = gin.Default()
	s.store = opts.Store
	s.log = opts.Log
	s.bot = opts.Bot

	s.GET("/status", s.pingHandler)

	s.GET("/send/:uuid", s.sendHandler)
}

func (s *service) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *service) sendHandler(c *gin.Context) {
	uuidString := c.Param("uuid")

	uid, err := uuid.Parse(uuidString)
	if err != nil {
		s.log.Error(
			"can't parse uuid",
			zap.String("uuid", uuidString),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't parse uuid",
		})

		return
	}

	text := c.Query("text")

	if text == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "text is empty",
		})

		return
	}

	record, err := s.store.GetByUUID(uid)
	if err != nil {
		s.log.Error(
			"can't get record",
			zap.Error(err),
			zap.String("uuid", uid.String()),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't get associated chat",
		})

		return
	}

	if err := s.bot.Send(record, text); err != nil {
		s.log.Error(
			"can't send text",
			zap.Error(err),
			zap.String("uuid", uid.String()),
			zap.String("text", text),
		)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "can't send text",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
