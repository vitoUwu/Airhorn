package airhorn

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/FedorLap2006/disgolf"
	"github.com/bwmarrin/discordgo"
)

var airhornBuffer = make([][]byte, 0)
var connections = make(map[string]string)

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
		guild, err := ctx.State.Guild(ctx.Interaction.GuildID)
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
			if errors.Is(err, discordgo.ErrStateNotFound) {
				ctx.Respond(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "You have to join in some voice channel",
					},
				})
			} else {
				log.Println(err)
				ctx.Respond(&discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "An error has occured while searching for your channel",
					},
				})
			}
			return
		}

		channelMention := fmt.Sprintf("<#%s>", voiceState.ChannelID)

		if _, alreadyConnected := connections[voiceState.ChannelID]; alreadyConnected {
			ctx.Respond(&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("I'm already playing in %s", channelMention),
				},
			})
			return
		}

		connections[voiceState.ChannelID] = voiceState.ChannelID
		defer delete(connections, voiceState.ChannelID)

		ctx.Respond(&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Playing airhorn in %s", channelMention),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})

		vc, err := ctx.ChannelVoiceJoin(voiceState.GuildID, voiceState.ChannelID, false, true)
		if err != nil {
			log.Print(err)
			content := fmt.Sprintf("An error has occured while joining in %s", channelMention)
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
	}),
}
