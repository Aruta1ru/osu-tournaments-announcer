package events

import (
	"discord-go/db"
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func NotifyNewForumposts(session *discordgo.Session) {
	forumpostIDs, err := db.GetNonNotifiedForumpostsID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(forumpostIDs) == 0 {
		return
	}

	notifiedChannelsID, err := db.GetNotifiedChannellsID()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	var embeds []*discordgo.MessageEmbed
	for _, forumpostID := range forumpostIDs {
		embedMessage, err := db.GetForumpostDataById(forumpostID)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		embeds = append(embeds, embedMessage)
	}

	for _, channelID := range notifiedChannelsID {
		idStr := strconv.FormatInt(channelID, 10)
		for _, embed := range embeds {
			_, err = session.ChannelMessageSendEmbed(idStr, embed)
			if err != nil {
				session.ChannelMessageSend(idStr, "Oops, something went wrong... :scream:")
				fmt.Println("Error:", err)
				return
			}
		}

	}

	err = db.SetForumpostsNotified(forumpostIDs)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

}
