package airhorn

import (
	"encoding/binary"
	"io"
	"log"
	"os"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

var airhornBuffer = make([][]byte, 0)
var connections = make(map[string]discordgo.Channel)

func LoadSound() (err error) {
	file, err := os.Open("airhorn.dca")
	if err != nil {
		return err
	}

	var opuslen int16

	for {
		err = binary.Read(file, binary.LittleEndian, &opuslen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err = file.Close()
			if err != nil {
				return err
			}

			return nil
		}

		if err != nil {
			return err
		}

		inBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &inBuf)
		if err != nil {
			return err
		}

		airhornBuffer = append(airhornBuffer, inBuf)
	}
}

var Data = &disgolf.Command{
	Name:        "airhorn",
	Description: "Plays airhorn sound in your voice channel",
	Handler: disgolf.HandlerFunc(func(ctx *disgolf.Ctx) {
		guild, err := ctx.Guild(ctx.Interaction.GuildID)
		if err != nil {
			ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command must be executed in guilds",
				},
			})
			return
		}

		voiceState, err := ctx.State.VoiceState(guild.ID, ctx.Interaction.Member.User.ID)

		if err != nil {
			ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You have to join in some voice channel",
				},
			})
			return
		}

		channel, err := ctx.Channel(voiceState.ChannelID)
		if err != nil {
			log.Print(err)
			ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "An error has occured while searching for your channel",
				},
			})
			return
		}

		_, alreadyConnected := connections[channel.ID]
		if alreadyConnected {
			ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "I'm already playing in your channel",
				},
			})
			return
		}

		connections[channel.ID] = *channel
		defer delete(connections, channel.ID)

		ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Joining...",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		vc, err := ctx.ChannelVoiceJoin(voiceState.GuildID, voiceState.ChannelID, false, true)
		if err != nil {
			log.Print(err)
			content := "An error has occured while joining in your voice channel"
			ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{Content: &content})
			return
		}

		time.Sleep(250 * time.Millisecond)

		vc.Speaking(true)
		for _, buff := range airhornBuffer {
			vc.OpusSend <- buff
		}

		vc.Speaking(false)
		time.Sleep(250 * time.Millisecond)
		vc.Disconnect()

		content := "Played Airhorn in your voice channel"
		ctx.InteractionResponseEdit(ctx.Interaction, &discordgo.WebhookEdit{Content: &content})
	}),
}
