package discord

import (
	"fmt"
	"mi-c2/internal/env"
	"mi-c2/internal/logging"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/openlyinc/pointy"
)

var (
	s   *discordgo.Session
	log logging.Logger

	registeredCommands []*discordgo.ApplicationCommand

	commands = []discordgo.ApplicationCommand{
		{
			Name:        "ark",
			Description: "Archive Team - Magnus Institute",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Name:        "get",
					Description: "Get Resources",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "clusters",
							Description: "List all connected clusters",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "cluster",
							Description: "Get status of a connected cluster",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:         "name",
									Description:  "Cluster Name",
									Type:         discordgo.ApplicationCommandOptionString,
									Required:     true,
									Autocomplete: true,
								},
							},
						},
						{
							Name:        "projects",
							Description: "List all available ArchiveTeam projects",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
						},
						{
							Name:        "project",
							Description: "Get information for a specific ArchiveTeam project",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:         "name",
									Description:  "Project Name",
									Type:         discordgo.ApplicationCommandOptionString,
									Required:     true,
									Autocomplete: true,
								},
							},
						},
					},
				},
				{
					Name:        "set",
					Description: "Set Configuration",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "cluster",
							Description: "Get status of a connected cluster",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:         "cluster",
									Description:  "Cluster Name",
									Type:         discordgo.ApplicationCommandOptionString,
									Required:     true,
									Autocomplete: true,
								},
								{
									Name:         "project",
									Description:  "Project Name",
									Type:         discordgo.ApplicationCommandOptionString,
									Required:     true,
									Autocomplete: true,
								},
								{
									Name:        "workers",
									Description: "Number Of Workers",
									Type:        discordgo.ApplicationCommandOptionInteger,
									MinValue:    pointy.Float64(1),
									MaxValue:    200,
									Required:    true,
								},
								{
									Name:        "concurrency",
									Description: "Concurrent Items per Worker",
									Type:        discordgo.ApplicationCommandOptionInteger,
									MinValue:    pointy.Float64(1),
									MaxValue:    20,
									Required:    false,
								},
							},
						},
					},
				},
				{
					Name:        "stop",
					Description: "Stop Workers",
					Type:        discordgo.ApplicationCommandOptionSubCommandGroup,
					Options: []*discordgo.ApplicationCommandOption{
						{
							Name:        "cluster",
							Description: "Get status of a connected cluster",
							Type:        discordgo.ApplicationCommandOptionSubCommand,
							Options: []*discordgo.ApplicationCommandOption{
								{
									Name:         "cluster",
									Description:  "Cluster Name",
									Type:         discordgo.ApplicationCommandOptionString,
									Required:     true,
									Autocomplete: true,
								},
							},
						},
					},
				},
			},
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ark": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			switch i.Type {

			// HANDLE COMMAND REQUESTS
			case discordgo.InteractionApplicationCommand:
				options := i.ApplicationCommandData().Options

				switch options[0].Name {
				case "get":
					switch options[0].Options[0].Name {

					// GET ALL ONLINE CLUSTERS
					case "clusters":
						getClustersHandler(i.Interaction)

					// GET SINGLE CLUSTER
					case "cluster":
						getClusterHandler(i.Interaction)

					case "projects":
						getProjectsHander(i.Interaction)
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Command \"%s\" is unhandled", options[0].Name),
						},
					})

				case "set":
					switch options[0].Options[0].Name {
					case "cluster":
						setClusterHandler(i.Interaction)
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Command \"%s\" is unhandled", options[0].Name),
						},
					})

				case "stop":
					switch options[0].Options[0].Name {
					case "cluster":
						stopClusterHandler(i.Interaction)
					}

					s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
						Type: discordgo.InteractionResponseChannelMessageWithSource,
						Data: &discordgo.InteractionResponseData{
							Content: fmt.Sprintf("Command \"%s\" is unhandled", options[0].Name),
						},
					})
				}

			case discordgo.InteractionApplicationCommandAutocomplete:
				data := i.ApplicationCommandData()

				// Root Command
				switch data.Options[0].Name {
				case "get":
					// Cluster Autocomplete
					switch data.Options[0].Options[0].Name {
					case "cluster":
						clusterListAutocompleteResponse(i.Interaction)

					// Project List Autocomplete
					case "project":
						projectListAutocompleteResponse(i.Interaction)
					}

				case "set":
					switch data.Options[0].Options[0].Name {
					case "cluster":
						switch {
						// Cluster List Autocomplete
						case data.Options[0].Options[0].Options[0].Focused:
							clusterListAutocompleteResponse(i.Interaction)
							// Project List Autocomplete
						case data.Options[0].Options[0].Options[1].Focused:
							projectListAutocompleteResponse(i.Interaction)
						}
					}

				case "stop":
					switch data.Options[0].Options[0].Name {
					case "cluster":
						switch {
						// Cluster List Autocomplete
						case data.Options[0].Options[0].Options[0].Focused:
							clusterListAutocompleteResponse(i.Interaction)
						}
					}
				}
			}
		},
	}
)

func Connect() {
	log = logging.Log.With().Str("module", "discord").Logger()
	log.Info().Msg("connecting to discord")

	var err error
	s, err = discordgo.New("Bot " + env.DISCORD_TOKEN)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to discord")
		os.Exit(1)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info().Msg("connected")
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			log.Info().Str("name", i.ApplicationCommandData().Name).Msg("processed cmd")
			h(s, i)
		}
	})

	cmdIDs := make(map[string]string, len(commands))

	registeredCommands = make([]*discordgo.ApplicationCommand, len(commands))

	for _, cmd := range commands {
		rcmd, err := s.ApplicationCommandCreate(env.DISCORD_APP_ID, env.DISCORD_GUILD_ID, &cmd)
		if err != nil {
			log.Error().Msgf("Cannot create slash command %q: %v", cmd.Name, err)
		}

		cmdIDs[rcmd.ID] = rcmd.Name
	}

	err = s.Open()

	if err != nil {
		log.Error().Err(err).Msg("unable to connect to discord")
		os.Exit(1)
	}

	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, env.DISCORD_GUILD_ID, &v)
		if err != nil {
			log.Error().Err(err).Msgf("Cannot create '%v' command: %v", v.Name, err)
			os.Exit(1)
		}
		log.Info().Str("cmd", v.Name).Msg("registered command")

		registeredCommands[i] = cmd
	}

}

func Cleanup() {
	log.Info().Msg("cleaning up commands")

	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, env.DISCORD_GUILD_ID, v.ID)
		if err != nil {
			if err != nil {
				log.Error().Err(err).Msgf("Cannot remove '%v' command: %v", v.Name, err)
				os.Exit(1)
			}
		}
	}

	s.Close()
}
