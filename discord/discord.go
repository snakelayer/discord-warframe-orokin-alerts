package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

type Discord struct {
	session *discordgo.Session

	warframeChannelIds []string
}

func New(token string) *Discord {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.WithError(err).Fatal("could not initialize discord")
	}

	err = session.Open()
	if err != nil {
		log.WithError(err).Fatal("error creating discord session")
	}

	return &Discord{
		session:            session,
		warframeChannelIds: []string{},
	}
}

func (discord *Discord) Initialize() {
	discord.initializeWarframeChannels()
}

func (discord *Discord) initializeWarframeChannels() {
	guilds, err := discord.session.UserGuilds(100, "", "")
	if err != nil {
		log.WithError(err).Fatal("could not retrieve server list")
	}

	for _, guild := range guilds {
		log.WithField("guild", guild).Info("identifying channels from guild")

		channels, err := discord.session.GuildChannels(guild.ID)
		if err != nil {
			log.WithError(err).Error("could not retrieve channel list")
			continue
		}

		for _, channel := range channels {
			log.WithField("channel", channel).Debug("channel data")

			// https://discordapp.com/developers/docs/resources/channel#channel-object-channel-types
			if channel.Type != 0 {
				continue
			}

			if isWarframeChannel(channel.Name) {
				log.WithField("channel", channel).Info("using warframe channel")
				discord.warframeChannelIds = append(discord.warframeChannelIds, channel.ID)
			}
		}
	}
}

func (discord *Discord) Broadcast(message string) {
	for _, channelId := range discord.warframeChannelIds {
		messageResponse, _ := discord.session.ChannelMessageSendTTS(channelId, "potato alert")

		if messageResponse != nil && messageResponse.ID != "" {
			discord.session.ChannelMessageEdit(channelId, messageResponse.ID, message)
		} else {
			discord.session.ChannelMessageSend(channelId, message)
		}
	}
}

func (discord *Discord) SetAlertStatus() {
	discord.session.UpdateStatus(0, "active alert!")
}

func (discord *Discord) ResetStatus() {
	discord.session.UpdateStatus(0, "")
}

func (discord *Discord) Close() {
	discord.session.Close()
}

func isWarframeChannel(name string) bool {
	if strings.Contains(strings.ToLower(name), "warframe") {
		return true
	} else if strings.Contains(strings.ToLower(name), "orokin") {
		return true
	} else if strings.Contains(strings.ToLower(name), "potato") {
		return true
	}

	return false
}
