package events

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Ready(session *discordgo.Session, event *discordgo.Ready) {
	session.UpdateListeningStatus("Aruta1ru")
	for {
		time.Sleep(time.Second * 5)
		NotifyNewForumposts(session)
		time.Sleep(time.Minute * 30)
	}
}
