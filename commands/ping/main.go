package ping

import (
	"fmt"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

var Data = &disgolf.Command{
	Name:        "ping",
	Description: "Check bot's latency",
	Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
		start := time.Now()

		ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		content := fmt.Sprintf("Latency: %s\nAPI: %s", time.Since(start).Round(time.Millisecond), ctx.HeartbeatLatency().Round(time.Millisecond))

		ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{Content: &content})
	}),
}
