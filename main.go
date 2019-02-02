package main

import (
  "fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"os/signal"
	"encoding/json"
	"io/ioutil"
	"regexp"
	"strconv"
)

var (
	// Expression that triggers the Minesweeper command.
	commandTriggerRegex = regexp.MustCompile(`^!minesweeper ?([0-9]+)? ?([0-9]+)? ?([0-9]+)?$`)
	// Simple array which maps proximity numbers to their respective spoilered emoji.
	// Simply calling numbers[cellValue] gets the necessary string to add.
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
	cfg Config
)

// Config is the config.json structure.
type Config struct {
	// Bot token to use.
	Token string `json:"token"`
	// Grid size, in the form of [width, height].
	DefaultGridSize [2]int `json:"defaultGridSize"`
	// Default number of mines to place on the grid.
	DefaultMines int `json:"defaultMines"`
}

func main() {
	// read config
	cfgfile, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Failed to read config file:", err)
	}
	// parse config
	err = json.Unmarshal(cfgfile, &cfg)
	if err != nil {
		fmt.Println("Failed to parse config file:", err)
	}
	// set default values, if unset
	if cfg.DefaultGridSize == [2]int{0, 0} {
		cfg.DefaultGridSize = [2]int{10,10}
	}
	if cfg.DefaultMines == 0 {
		cells := cfg.DefaultGridSize[0] * cfg.DefaultGridSize[1]
		cfg.DefaultMines = cells / 10
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

// sendGrid is the discordgo handler function that actually generates and sends the Minesweeper grids.
func sendGrid(bot *discordgo.Session, evt *discordgo.MessageCreate) {
	// test the message
	groups := commandTriggerRegex.FindStringSubmatch(evt.Message.Content)
	// if this returns nil, we don't have a match
	if groups == nil {
		return
	}
	// get default x, y, and mines values
	x := cfg.DefaultGridSize[0]
	y := cfg.DefaultGridSize[1]
	mines := cfg.DefaultMines
	// if all 3 parameters are set, then parse the user's parameters to ints and use them to generate the grid
	if groups[1] != "" && groups[2] != "" && groups[3] != "" {
		var err error
		x, err = strconv.Atoi(groups[1])
		if err != nil {
			fmt.Println("Error converting X string: ",err)
		}
		y, err = strconv.Atoi(groups[2])
		if err != nil {
			fmt.Println("Error converting Y string: ",err)
		}
		mines, err = strconv.Atoi(groups[3])
		if err != nil {
			fmt.Println("Error converting mines string: ",err)
		}
		// if this doesn't work and the 3 groups aren't equal to "", then this is not a valid input
	} else if groups[1] != groups[2] || groups[2] != groups[3] || groups[3] != groups[1] {
		fmt.Println("invalid input")
		SendErrorMessage(bot, evt, "invalid input")
	}
	// validate the user's input
	if x > 20 || y > 20 {
		SendErrorMessage(bot, evt, "that board's too big!")
		return
	}
	if x < 1 || y < 1 {
		SendErrorMessage(bot, evt, "you can't make a board with no cells!")
		return
	}
	if mines >= (x * y) {
		SendErrorMessage(bot, evt, "you can't make a board with as many mines as there are cells!")
		return
	}
	if mines == 0 {
		SendErrorMessage(bot, evt, "you can't make a board with no mines!")
		return
	}
	// generate the grid
	fmt.Println("making a new board for",evt.Message.Author.Username+"#"+evt.Message.Author.Discriminator)
	grid := NewMSGrid(x, y)
	grid.Populate(mines)
	grid.updateMineCount()
	dgrid := DiscordGrid(grid)
	// if it's too big, we can't send it
	if len(dgrid) > 2000 {
		SendErrorMessage(bot, evt, "that board's too big!")
	}
	_, err := bot.ChannelMessageSend(evt.Message.ChannelID,dgrid)
	if err != nil {
		fmt.Println("Error sending board:",err)
	}
}

// SendErrorMessage @s the user in the event and tells them what went wrong with the command.
func SendErrorMessage(bot *discordgo.Session, evt *discordgo.MessageCreate, msg string) {
	_, err := bot.ChannelMessageSend(evt.Message.ChannelID,fmt.Sprintf("<@%s>, %s",evt.Message.Author.ID, msg))
	if err != nil {
		fmt.Println("Error sending error message:",err)
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
