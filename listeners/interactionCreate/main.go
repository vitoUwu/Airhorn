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

	guild, err := s.Guild(i.Interaction.GuildID)
	if err == nil {
		log.Printf("[GUILD] %s command executed at %s", commandName, guild.Name)
	} else if i.Interaction.User != nil {
		user, err := s.User(i.Interaction.User.ID)
		if err == nil {
			log.Printf("[DM] %s command executed at %s DM", commandName, user.Username)
		} else {
			log.Printf("[Unknown] %s command executed", commandName)
		}
	} else {
		log.Printf("[Unknown] %s command executed", commandName)
	}
}
