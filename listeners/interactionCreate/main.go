package interactionCreate

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

var Handler = func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.Interaction.ApplicationCommandData().Name

	guild, err := s.State.Guild(i.Interaction.GuildID)
	if err == nil {
		log.Printf("[GUILD] %s command executed at %s", commandName, guild.Name)
	} else if user := i.Interaction.User; user != nil {
		log.Printf("[DM] %s command executed at %s DM", commandName, user.Username)
	} else {
		log.Printf("[Unknown] %s command executed", commandName)
	}
}
