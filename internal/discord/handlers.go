package discord

import (
	"fmt"
	"mi-c2/internal/controller"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getClustersHandler(interaction *discordgo.Interaction) {
	clusters := controller.GetClusters()
	if len(clusters) == 0 {
		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No Clusters Online",
			},
		})
		return
	} else {
		res := []string{}
		for _, cluster := range clusters {
			res = append(res, fmt.Sprintf("%s, owners %+v, current project: %s (%d/%d workers @ %d concurrency), last modified by %s", cluster.Name, cluster.Owners, cluster.Project, cluster.ActiveWorkers, cluster.Workers, cluster.Concurrency, cluster.Cause))
		}

		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: strings.Join(res, "\n"),
			},
		})

		return
	}
}

func getClusterHandler(interaction *discordgo.Interaction) {
	options := interaction.ApplicationCommandData().Options
	clusterName := options[0].Options[0].Options[0].StringValue()

	cluster, err := controller.GetCluster(clusterName)
	if err != nil {
		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Cluster %s Is Offline", clusterName),
			},
		})

		return
	} else {
		s.InteractionRespond(interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("%s, owners %+v, current project: %s (%d/%d workers @ %d concurrency), last modified by %s", cluster.Name, cluster.Owners, cluster.Project, cluster.ActiveWorkers, cluster.Workers, cluster.Concurrency, cluster.Cause),
			},
		})

		return
	}
}

func getProjectsHander(interaction *discordgo.Interaction) {
	projects := controller.GetProjects()
	res := []string{}

	res = append(res, fmt.Sprintf("DEFAULT PROJECT: **%s**\n", projects.Default))

	for _, project := range projects.Projects {
		res = append(res, fmt.Sprintf("**%s** (%s): %s", project.Title, project.Name, project.Description))
	}

	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: strings.Join(res, "\n"),
		},
	})
}

func setClusterHandler(interaction *discordgo.Interaction) {
	options := interaction.ApplicationCommandData().Options
	cluster := options[0].Options[0].Options[0].StringValue()
	project := options[0].Options[0].Options[1].StringValue()
	workers := options[0].Options[0].Options[2].IntValue()
	concurrency := options[0].Options[0].Options[3].IntValue()

	if concurrency == 0 {
		concurrency = 4
	}

	controller.RequestClusterChange(interaction.Member.User.Username, cluster, project, int(workers), int(concurrency))

	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Setting **%s** to `%s` with %d workers @ 5 concurrent", cluster, project, workers),
		},
	})
}

func stopClusterHandler(interaction *discordgo.Interaction) {
	options := interaction.ApplicationCommandData().Options
	cluster := options[0].Options[0].Options[0].StringValue()

	controller.RequestClusterChange(interaction.Member.User.Username, cluster, "none", 0, 0)

	s.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Stopping all workers on **%s**", cluster),
		},
	})
}
