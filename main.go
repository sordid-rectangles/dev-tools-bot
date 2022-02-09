package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const Version = "v0.0.1-alpha"

var dg *discordgo.Session
var TOKEN string
var GUILDID string

func init() {
	// Print out a fancy logo!
	fmt.Printf(`Science Defender! %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	//Load Token from env (simulated with godotenv)
	TOKEN = os.Getenv("BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading token from env")
		os.Exit(1)
	}

	GUILDID = os.Getenv("GUILD_ID")
	if GUILDID == "" {
		log.Println("No GuildID specified in env")
		GUILDID = "" //this effectively specifies command registration as global
	}
}

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "whoami",
			Description: "Returns information about the user who called the command",
		},
		{
			Name:        "chan",
			Description: "Returns information about the channel the command is called in",
		},
		{
			Name:        "serv",
			Description: "Returns information about the server the command is called in",
		},
		{
			Name:        "user-info",
			Description: "Returns information about the specified user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-select",
					Description: "select a user to register",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whoami": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.Member.User.ID
			nick := i.Member.Nick
			name := i.Member.User.Username
			join := i.Member.JoinedAt

			content := `Your discord user id is: %s
			Your Nick is: %s
			Your Username is: %s
			You joined this server: %s`
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, id, nick, name, join),
				},
			})
		},
		"chan": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.Message.ChannelID

			content := `Channel discord id is: %s`
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, id),
				},
			})
		},
		"serv": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.GuildID
			locale := i.GuildLocale

			content := `Guild discord id is: %s
			Guild locale is: %s`
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, id, locale),
				},
			})
		},
		"user-info": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(s)
			id := user.ID
			name := user.Username
			str := user.String()

			content := `User discord id is: %s
			User Username is: %s
			User str is: %s`
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, id, name, str),
				},
			})
		},
	}
)

func init() {
	var err error
	dg, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating discordgo session!")
		os.Exit(1)
	}
}

func main() {
	var err error
	//Configure discordgo session bot
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//Register Bot Intents with Discord
	//worth noting MakeIntent is a no-op, but I want it there for doing something with pointers later
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GUILDID, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)

	defer dg.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")

	// Exit Normally.
	//exit
}
