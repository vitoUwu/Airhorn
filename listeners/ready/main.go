package ready

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var Handler = func(s *discordgo.Session, r *discordgo.Ready) {
	log.Printf("Ready as %s ", s.State.User.Username)
	guilds := len(s.State.Guilds)
	s.UpdateGameStatus(0, fmt.Sprintf("in %d servers", guilds))
}
