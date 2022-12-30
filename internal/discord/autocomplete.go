package discord

import (
	"mi-c2/internal/controller"

	"github.com/bwmarrin/discordgo"
)

func clusterListAutocompleteResponse(i *discordgo.Interaction) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	clusters := controller.GetClusters()
	for _, cluster := range clusters {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  cluster.Name,
			Value: cluster.Name,
		})
	}

	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}

func projectListAutocompleteResponse(i *discordgo.Interaction) {
	var choices []*discordgo.ApplicationCommandOptionChoice
	projects := controller.GetProjects()
	for _, project := range projects.Projects {
		choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
			Name:  project.Title,
			Value: project.Name,
		})
	}

	s.InteractionRespond(i, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
}
