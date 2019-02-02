package main

import (
  "fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"encoding/json"
	"io/ioutil"
)

var (
	numbers = [...]string{
		"||:zero:||",
		"||:one:||",
		"||:two:||",
		"||:three:||",
		"||:four:||",
		"||:five:||",
		"||:six:||",
		"||:seven:||",
		"||:eight:||",
		"||:bomb:||",
	}
)

type Config struct {
	Token string `json:"token"`
}

func main() {
	// read config
	cfgfile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Failed to read config file:", err)
	}
	// parse config
	var cfg Config
	err = json.Unmarshal(cfgfile, &cfg)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
	}
	// start bot
	bot, err := discordgo.New("Bot "+cfg.Token)
	if err != nil {
		fmt.Println("Failed to start bot:", err)
	}
	// add bot handler
	bot.AddHandler(sendGrid)
	// start listening for messages
	bot.Open()
	defer bot.Close()
	// notify user, wait for interrupt to quit the bot
	fmt.Println("ready (Ctrl-c to quit)")
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt, os.Kill)
	<-sch
	fmt.Println("exiting")
}

func sendGrid(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	if evt.Message.Content != "!minesweeper" {
		return
	}
	fmt.Println("making a new board for",evt.Message.Author.Username+"#"+evt.Message.Author.Discriminator)
	grid := NewMSGrid(10, 10)
	grid.Populate(20)
	grid.updateMineCount()
	_, err := bot.ChannelMessageSend(evt.Message.ChannelID,DiscordGrid(grid))
	if err != nil {
		fmt.Println(err)
	}
}

func DiscordGrid(grid MSGrid) (str string) {
	for _, row := range grid {
		for _, col := range row {
			str += fmt.Sprint(numbers[col])
		}
		str += fmt.Sprint("\n")
	}
	return str
}
