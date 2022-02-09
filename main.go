package main

import (
	"bytes"
	"fmt"
	"messfar-discord/domain"
	"messfar-discord/repo"
	"messfar-discord/util"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type MessageHandler struct {
	domain.FaceService
}

// Variables used for command line parameters
var (
	Token string
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Sprintf("fatal error config. %+v\n", err))
	}
	viper.AutomaticEnv()

	Token = viper.GetString("DISCORD_BOT_TOKEN")

	faceService := repo.NewFaceService(viper.GetString("FACE_SERVICE"))

	util.Init()

	messageHandler := MessageHandler{FaceService: faceService}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageHandler.ReceiveMessage)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func (mh *MessageHandler) ReceiveMessage(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if mh.IsOneImageMessage(m) {
		mh.HandleAttachments(s, m)
		return
	} else if mh.IsTextMessage(m) {
		mh.HandleMessage(s, m)
		return
	}
}

func (mh *MessageHandler) IsTextMessage(m *discordgo.MessageCreate) bool {
	return m.Content != ""
}

func (mh *MessageHandler) IsOneImageMessage(m *discordgo.MessageCreate) bool {
	return len(m.Attachments) == 1 && util.IsImage(m.Attachments[0].URL)
}

func (mh *MessageHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	switch m.Content {
	case "許願":
		getRandomResponse, err := mh.FaceService.GetRandom(1)
		if err != nil {
			log.Errorf("get random failed. error: %+v", err)
			return
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s`\n點我看資料: https://messfar.com/?ID=%s", getRandomResponse[0].Name, getRandomResponse[0].ID))
		s.ChannelMessageSend(m.ChannelID, getRandomResponse[0].Preview)
	case "我心愛的女孩":
		s.ChannelMessageSend(m.ChannelID, "請至髒沙發頁面觀看: https://messfar.com")
	}

}

func (mh *MessageHandler) HandleAttachments(s *discordgo.Session, m *discordgo.MessageCreate) {
	imageBytes, err := util.DownloadImage(m.Attachments[0].URL)
	if err != nil {
		return
	}

	// 如果圖片過大就縮小
	width, height, err := util.GetSize(bytes.NewBuffer(imageBytes))
	if err != nil {
		log.Errorf("get size failed. error: %+v", err)
		return
	}

	log.Info("image size: ", width, height)

	resizeImageBuffer := bytes.NewBuffer(imageBytes)
	if width > 600 || height > 600 {
		resizeImageBuffer, err = util.ImageResizeByBuffer(bytes.NewBuffer(imageBytes), 600)
		if err != nil {
			log.Errorf("resize image failed. error: %+v", err)
			return
		}
		// resizeImageBuffer, err = util.ImageProcessing(imageBytes, 100)
		// if err != nil {
		// 	log.Errorf("resize image failed. error: %+v", err)
		// 	return
		// }
	}

	postSearchResponse, err := mh.FaceService.PostSearch(resizeImageBuffer.Bytes())
	if err != nil {
		log.Errorf("post search failed. error: %+v", err)
		return
	}

	for _, v := range postSearchResponse {
		s.ChannelMessageSend(m.ChannelID,
			fmt.Sprintf("我猜可能是`%s`\n相似度: %f%%\n點我看資料: https://messfar.com/?ID=%s", v.Name, v.RecognitionPercentage, v.ID),
		)
		s.ChannelMessageSend(m.ChannelID, v.Preview)
	}
}
