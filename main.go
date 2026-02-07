package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	disgo "github.com/bwmarrin/discordgo"
	env "github.com/joho/godotenv"
)

type Reminder struct {
	Name	string
	Id	int
	TimeCreated	time.Time
	TimeOfReminder	time.Time
}

var _master_id = 0
var _reminders = make([]Reminder, 0)

func init() {

	// Load environment variables from the .env file
	err := env.Load()
	if err != nil {
		log.Fatal("---ERROR: problem loading .env file")
	}
}

func main() {

	/* --- TEST --- */
	_reminders = append(_reminders, Reminder{"Reminder #1", _master_id, time.Now(), time.Now()})
	_master_id = _master_id + 1
	_reminders = append(_reminders, Reminder{"Reminder #2", _master_id, time.Now(), time.Now()})
	_master_id = _master_id + 1

	// Get the bot token from environment variable
	Token := os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		fmt.Println("---ERROR: No token provided. Set the DISCORD_TOKEN environment variable.")
		return
	}

	// Create a new Discord session using the provided bot token.
	dg, err := disgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("---ERROR: could not create Discord session,", err)
		return
	}

	// Register the messageCreate function as a callback for the MessageCreate event.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("---ERROR: problem opening connection,", err)
		return
	}
	defer dg.Close()

	fmt.Println("* Bot is now running. Press CTRL+C to exit.")
	select {}
}

// This function will be called every time a new message is created in the Discord server.
func messageCreate(s *disgo.Session, m *disgo.MessageCreate) {
	// Ignore messages from the bot itself.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Connectivity Check
	if m.Content == "r-ping" {
		s.ChannelMessageSend(m.ChannelID, "Connected")
	}

	if m.Content == "r-help" {
		usage := `## Remi Bot Usage
		*r-help*	: Display this message
		*r-ping*	: Check bot connectivity`
		s.ChannelMessageSend(m.ChannelID, usage)
	}

	if m.Content == "r-embed" {
		embed := &disgo.MessageEmbed{
			Title:	"Sample Embed",
			Description:	"This is an embedded message!",
			Color:	0x00ff00, // Green
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if m.Content == "r-printAll" {
		reminderString := "* "
		for i := 0; i < len(_reminders); i++ {
			if i != len(_reminders)-1 {
				reminderString = reminderString + _reminders[i].Name + ": " + _reminders[i].TimeOfReminder.Truncate(time.Second).String() + "\n* "
			} else {
				reminderString = reminderString + _reminders[i].Name + ": " + _reminders[i].TimeOfReminder.Truncate(time.Second).String() + "\n"
			}
		}
		embed := &disgo.MessageEmbed{
			Title: " Scheduled Reminders",
			Description: reminderString,
			Color: 0x00ff00, // Green
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}

	if strings.Contains(m.Content, "r-schedule") {
		splitMsg := strings.Split(m.Content, "-")
		var newReminder Reminder
		newReminder.Id = _master_id
		_master_id = _master_id+1
		newReminder.TimeCreated = time.Now()
		newReminder.TimeOfReminder = time.Now()

		for i:= 1; i < len(splitMsg); i++ {

			//parse inputs
			fmt.Println(splitMsg[i])

			// -d 
			// Date

			// -t
			// Time

			// -n
			// Name
			if splitMsg[i][0:2] == "n " {
				name, found := strings.CutPrefix(splitMsg[i], "n ")
				if found {
					newReminder.Name = name
				}
			}
		}

		_reminders = append(_reminders, newReminder)

		s.ChannelMessageSend(m.ChannelID, "Added new reminder for "+newReminder.Name)
	}
} 
