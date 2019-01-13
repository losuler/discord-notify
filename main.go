package main

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
	"log"
	"flag"

	"gopkg.in/bwmarrin/discordgo.v0"
	"gopkg.in/go-telegram-bot-api/telegram-bot-api.v4"
)

var verbosePtr bool

type Conf struct {
	DiscordToken string `json:"discordToken"`
	TelegramToken string `json:"telegramToken"`
	TelegramChatID int64 `json:"telegramChatID"`
}

func createLogin() {
	var email, pass, token string

	// TODO: Check if works with passwords with spaces.
	// https://stackoverflow.com/q/20895552
	fmt.Println("Enter email:")
	fmt.Scanln(&email)
	// TODO: Hide password input.
	// https://stackoverflow.com/q/2137357
	fmt.Println("Enter password:")
	fmt.Scanln(&pass)

	dg, err := discordgo.New(email, pass)
	if err != nil {
		fmt.Println("Error creating Discord session using email/pass,", err)
		os.Exit(1)
	}
	// Obtaining authentication token for Discord account.
	token = dg.Token
	input := Conf{DiscordToken: token}
	
	data, err := json.MarshalIndent(input, "", "\t")
	if err != nil {
		fmt.Println("Error encoding JSON,", err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("conf.json", data, 0600)
	if err != nil {
		fmt.Println("Error writing JSON file,", err)
		os.Exit(1)
	}
}

func retrieveFromJSON() (string, string, int64) {
	var conf Conf

	raw, err := ioutil.ReadFile("conf.json")
	if err != nil {
		fmt.Println("Error reading JSON file,", err)
		os.Exit(1)
	}

	err = json.Unmarshal(raw, &conf)
	if err != nil {
		fmt.Println("Error decoding JSON,", err)
		os.Exit(1)
	}

	log.Println(conf.DiscordToken, conf.TelegramToken, conf.TelegramChatID)
	return conf.DiscordToken, conf.TelegramToken, conf.TelegramChatID
}

func processMessages(discordToken string) {
	dg, err := discordgo.New(discordToken)
	if err != nil {
		fmt.Println("Error creating Discord session using token,", err)
		os.Exit(1)
	}

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection,", err)
		os.Exit(1)
	}

	// Register the directRecieve func as a callback for MessageCreate events.
	dg.AddHandler(directRecieve)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	// Channel to recieve signal notifications.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sigChan

	dg.Close()
}

func directRecieve(s *discordgo.Session, m *discordgo.MessageCreate) {
	var telegramToken string
	var telegramChatID int64
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("Error", err)
		os.Exit(1)
	}
	if m.Author.ID == s.State.User.ID {
		return
	}	
	if channel.Type == discordgo.ChannelTypeDM {
		log.Println(m.Author.Username)
		log.Println(m.Content)
		log.Println(time.Now().Format(time.RFC850))
		_, telegramToken, telegramChatID = retrieveFromJSON()
		telegramNotify(m.Author.Username, telegramToken, telegramChatID)
	} else {
		return
	}
}

func telegramNotify(message, telegramToken string, telegramChatID int64) {
	bot, err := tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		fmt.Println("Error authenticating bot,", err)
		os.Exit(1)
	}
	request := tgbotapi.NewMessage(telegramChatID, message)
	bot.Send(request)
}

func main() {
	var discordToken string
	
	verbosePtr := flag.Bool("verbose", false, "Print verbose log output")
	flag.Parse()

	log.SetOutput(ioutil.Discard)
	if *verbosePtr {
		log.SetOutput(os.Stderr)
	}
	
	_, err := os.Stat("conf.json")
	if os.IsNotExist(err) {
		fmt.Println("File doesn't exist, creating file...")
		createLogin()
		discordToken, _, _ = retrieveFromJSON()
	} else if !os.IsNotExist(err) {
		fmt.Println("File exists, accessing...")
		discordToken, _, _ = retrieveFromJSON()
	} else {
		fmt.Println("File stat error,", err)
	}

	processMessages(discordToken)
}
