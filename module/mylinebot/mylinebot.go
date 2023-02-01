package mylinebot

import (
	"context"
	"github.com/line/line-bot-sdk-go/linebot"
	"log"
	"openai-line-bot/clients"
	gpt3 "openai-line-bot/clients/gp3"
	"strings"
)

func LineBotTemplate(events []*linebot.Event) {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				_, err := clients.MyLineBot.GetMessageQuota().Do()
				if err != nil {
					log.Println("Quota err:", err)
				}
				log.Println("receive message :: ", message.Text)

				isGroup := len(event.Source.GroupID) > 0
				if !strings.Contains(message.Text, "@bot") && isGroup {
					log.Println("dont container @bot , just return")
					continue
				}
				// 將訊息丟給OpenA
				isImage := strings.Contains(message.Text, "圖") ||
					strings.Contains(message.Text, "@botimg") ||
					strings.Contains(message.Text, "照片")

				txt := strings.Replace(message.Text, "@botimg", "", 1)
				txt = strings.Replace(message.Text, "@bot", "", 1)
				if isImage {
					response, err := requestImageFromOpenAI(txt)

					if err != nil {
						if _, err = clients.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
							log.Print(err)
						}
					}
					log.Println("send image ...")
					// 回覆OpenAI回應的訊息
					if _, err = clients.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(response, response)).Do(); err != nil {
						log.Print(err)
					}
					continue
				}
				openAIresp := requestOpenAI(txt)

				log.Println("send message :: ", openAIresp[0:5])

				if _, err = clients.MyLineBot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(openAIresp)).Do(); err != nil {
					log.Print(err)
				}

			}
		}
	}
}

func requestOpenAI(line_Message string) string {
	ctx := context.Background()

	resp, err := clients.MyOpenAI.CompletionWithEngine(ctx, "text-davinci-003", gpt3.CompletionRequest{
		Prompt:    []string{line_Message},
		MaxTokens: gpt3.IntPtr(512),
		//Stop:      []string{"."},
		//Echo:      false,
	})

	if err != nil {
		log.Fatalln(err)
	}
	return resp.Choices[0].Text

}

func requestImageFromOpenAI(line_Message string) (string, error) {
	ctx := context.Background()

	request := gpt3.ImageRequest{
		Prompt: line_Message,
		Number: 1,
		Size:   "512x512",
	}
	resp, err := clients.MyOpenAI.Image(ctx, request)

	if err != nil {
		log.Fatalln(err)
	}
	return resp.Data[0].Url, err

}
