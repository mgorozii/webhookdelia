package server

import (
	"bytes"
	"html/template"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mgorozii/webhookdelia/internal/services/bot"
	"github.com/mgorozii/webhookdelia/internal/services/store"
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
	s.POST("/send/:uuid", s.sendHandler)
}

func (s *service) pingHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

// nolint:funlen // todo
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

	data, err := newRequestData(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	t, err := template.New("").Parse(data.Text)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	buf := bytes.NewBuffer(nil)

	if err := t.Execute(buf, data.Body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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

	sendBuffer := make([]byte, 4096)

	for {
		n, err := buf.Read(sendBuffer)

		if err != nil && err != io.EOF {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})

			return
		}

		if n == 0 {
			break
		}

		if err := s.bot.Send(record, string(sendBuffer[:n]), data.ParseMode); err != nil {
			s.log.Error(
				"can't send text",
				zap.Error(err),
				zap.String("uuid", uid.String()),
				zap.String("text", string(sendBuffer[:n])),
			)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "can't send text",
			})

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
