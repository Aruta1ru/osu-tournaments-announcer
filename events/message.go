package events

import (
	"discord-go/db"
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "$notify" {
		guildData, err := s.Guild(m.GuildID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Oops, something went wrong... :scream:")
			return
		}

		if m.Author.ID != guildData.OwnerID {
			s.ChannelMessageSend(m.ChannelID, "Sorry, this command available only for server owner")
			return
		}

		serverID, _ := strconv.ParseInt(m.GuildID, 10, 0)
		channelID, _ := strconv.ParseInt(m.ChannelID, 10, 0)

		message, err := db.NotifyServerChannel(serverID, channelID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Oops, something went wrong... :scream:")
			return
		}

		s.ChannelMessageSend(m.ChannelID, message)
		return
	}

	messageStrings := strings.Split(m.Content, " ")
	if messageStrings[0] == "$gp" {
		if len(messageStrings) == 1 {
			s.ChannelMessageSend(m.ChannelID, "No forumpost ID found.\nUse '$gp <forumpost ID>'")
			return
		}
		forumpostID, err := strconv.Atoi(messageStrings[1])
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Please, provide a number to find forumpost data")
			return
		}
		embedMessage, err := db.GetForumpostDataById(forumpostID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Oops, something went wrong... :scream:")
			fmt.Println("Error:", err)
			return
		}
		if embedMessage == nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Forumpost with ID: %d not found", forumpostID))
			return
		}
		_, err = s.ChannelMessageSendEmbed(m.ChannelID, embedMessage)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Oops, something went wrong... :scream:")
			fmt.Println("Error:", err)
			return
		}
		return
	}

}
