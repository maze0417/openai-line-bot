package mybot

import (
	"fmt"
	"openai-line-bot/clients"
	"openai-line-bot/module/mylinebot"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v7/linebot"
)

// var bot *linebot.Client
func NewStart(r *gin.Context) {
	events, err := clients.MyLineBot.ParseRequest(r.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			// w.WriteHeader(400)
			fmt.Println(400)
		} else {
			// w.WriteHeader(500)
			fmt.Println(500)
		}
		return
	}

	mylinebot.LineBotTemplate(events)
}
