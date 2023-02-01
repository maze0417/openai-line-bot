package cmd

import (
	"fmt"
	"log"
	"openai-line-bot/clients"
	"openai-line-bot/controller/mybot"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// server command
var serverCmd = &cobra.Command{
	Use: "server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

// add command
func init() {
	rootCmd.AddCommand(serverCmd)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .clients file")

	}
}

func start() {
	fmt.Println("start function")

	clients.LineConn()
	clients.Gpt3Conn()

	ginServer := gin.New()
	ginServer.SetTrustedProxies(nil)
	ginServer.POST("/callback", mybot.NewStart)
	ginServer.GET("/", func(r *gin.Context) {
		r.JSON(200, gin.H{"message": "ai bot ready", "code": 0})
	})
	ginServer.Run(":8833")
}
