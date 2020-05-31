package server

import (
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/juju/errors"
	"gopkg.in/tucnak/telebot.v2"
)

func newRequestData(c *gin.Context) (requestData, error) {
	data := requestData{
		Body: map[string]interface{}{},
	}

	if err := c.ShouldBind(&data.Body); err != nil {
		return data, err
	}

	getKey := func(key string) string {
		value := c.Query(key)
		if value == "" {
			value, _ = data.Body[key].(string)
			delete(data.Body, key)
		}

		return value
	}

	data.Text = getKey("text")

	data.ParseMode = getKey("parse_mode")

	if err := data.Validate(); err != nil {
		return data, errors.Trace(err)
	}

	return data, nil
}

type requestData struct {
	Text      string
	ParseMode telebot.ParseMode
	Body      map[string]interface{}
}

func (rd requestData) Validate() error {
	return validation.ValidateStruct(&rd,
		validation.Field(&rd.Text, validation.Required, validation.Length(1, 0)),
		validation.Field(
			&rd.ParseMode,
			validation.In(
				telebot.ModeDefault,
				telebot.ModeMarkdown,
				telebot.ModeMarkdownV2,
				telebot.ModeHTML,
			),
		),
	)
}
