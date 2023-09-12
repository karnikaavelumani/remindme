package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/api/cmdroute"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/joho/godotenv"
)

var commands = []api.CreateCommandData{{
	Name:        "remindme",
	Description: "Set a reminder",
	Options: discord.CommandOptions{
		&discord.SubcommandOption{
			OptionName:  "remind",
			Description: "Set a reminder",
		},
		&discord.SubcommandOption{
			OptionName:  "recur",
			Description: "Set a recurring reminder",
		},
	},
}}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: failed to load .env:", err)
	}

	r := cmdroute.NewRouter()
	r.AddFunc("remindme", func(ctx context.Context, data cmdroute.CommandData) *api.InteractionResponseData {
		return &api.InteractionResponseData{Content: option.NewNullableString("Pong!")}
	})

	s := state.New("Bot " + os.Getenv("BOT_TOKEN"))
	s.AddInteractionHandler(r)
	s.AddIntents(gateway.IntentGuilds)
	s.AddHandler(func(e *gateway.InteractionCreateEvent) {
		var resp api.InteractionResponse

		switch data := e.Data.(type) {
		case *discord.CommandInteraction:
			jsonContent, err := json.MarshalIndent(data, "", "  ")
			if err != nil {
				log.Println("cannot marshal data:", err)
			}
			log.Println(string(jsonContent))
			resp = api.InteractionResponse{
				Type: api.MessageInteractionWithSource,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString("Unknown command: " + data.Name),
				},
			}

		case discord.ComponentInteraction:
			resp = api.InteractionResponse{
				Type: api.UpdateMessage,
				Data: &api.InteractionResponseData{
					Content: option.NewNullableString("Custom ID: " + string(data.ID())),
				},
			}
		default:
			log.Printf("unknown interaction type %T", e.Data)
			return
		}

		if err := s.RespondInteraction(e.ID, e.Token, resp); err != nil {
			log.Println("failed to send interaction callback:", err)
		}
	})

	if err := cmdroute.OverwriteCommands(s, commands); err != nil {
		log.Fatalln("cannot update commands:", err)
	}

	if err := s.Connect(context.TODO()); err != nil {
		log.Println("cannot connect:", err)
	}
}
