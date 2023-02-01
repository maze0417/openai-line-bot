package mylinebot

import (
	"context"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"openai-line-bot/env"
	"strings"
)

func LineBotTemplate(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				// quota, err := env.MyLineBot.GetMessageQuota().Do()
				_, err := env.MyLineBot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				log.Println("receive message :: ", message.Text)

				isGroup := len(event.Source.GroupID) > 0
				if !strings.Contains(message.Text, "@bot") && isGroup {
					log.Println("dont container @bot , just return")
					continue
				}

				// 將訊息丟給OpenAI
				openAIresp := requestOpenAI(message.Text)

				// 回覆同樣的訊息
				//if _, err = env.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
				//	log.Print(err)
				//}

				log.Println("send message :: ", openAIresp[0:5])
				// 回覆OpenAI回應的訊息
				if _, err = env.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(openAIresp)).Do(); err != nil {
					log.Print(err)
				}
				//case *linebot.StickerMessage:
				//	var kw string
				//	for _, k := range message.Keywords {
				//		kw = kw + "," + k
				//	}
				//
				//	outStickerResult := fmt.Sprintf("收到貼圖訊息: %s, pkg: %s kw: %s", message.StickerID, message.PackageID, kw)
				//	if _, err := env.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outStickerResult)).Do(); err != nil {
				//		log.Print(err)
				//	}
			}
		}
	}
}

func requestOpenAI(line_Message string) string {
	ctx := context.Background()
	resp, err := env.MyOpenAI.CompletionWithEngine(ctx, "text-davinci-003", gpt3.CompletionRequest{
		Prompt:    []string{line_Message},
		MaxTokens: gpt3.IntPtr(512),
		//Stop:      []string{"."},
		//Echo:      false,
	})

	if err != nil {
		log.Fatalln(err)
	}
	return resp.Choices[0].Text
	//log.Println("get raw from ai:: ", resp.Choices[0].Text)
	//res1 := strings.Split(, "\r\n")
	//if len(res1) > 1 {
	//	return res1[1]
	//}
	//
	//return res1[0]
}
