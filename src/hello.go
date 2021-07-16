package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var token = "ODY0ODU5OTM1MTY2NDMxMjUz.YO7llA.xZ9IxKr8PfxZQREubIVqnDecGP4"

func main() {

	dg, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

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

func ready(s *discordgo.Session, event *discordgo.Ready) {

	// Set the playing status.
	s.UpdateGameStatus(0, "Albion Online")
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.HasPrefix(m.Content, "!register") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		usr, err := s.User(m.Author.ID)
		if err != nil {
			fmt.Println("error on email get")
			return
		}
		s.ChannelMessageSend(c.ID, "<@"+m.Author.ID+"> Your name has been registered. If your character is in the Alliance, You will automatically be registered. this can take up to 5 minutes, Please be patience.")

		//Need to go into albion online api and look for ig name and guild name.
		s.GuildMemberNickname(c.GuildID, usr.ID, "")

		var inGame = strings.Split(m.Content, " ")

		response, err := http.Get("https://gameinfo.albiononline.com/api/gameinfo/search?q=" + inGame[1])

		responseData, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Print(err.Error())
			os.Exit(1)
		}

		fmt.Println(string(responseData))

	}

	if strings.HasPrefix(m.Content, "!unregister") {

		// Find the channel that the message came from.
		c, err := s.State.Channel(m.ChannelID)
		if err != nil {
			// Could not find channel.
			return
		}

		usr, err := s.User(m.Author.ID)
		if err != nil {
			fmt.Println("error on email get")
			return
		}

		s.GuildMemberNickname(c.GuildID, usr.ID, "Atns Child")

	}
}
