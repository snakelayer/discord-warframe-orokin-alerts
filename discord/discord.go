package discord

import (
	"errors"
	"strings"

	"github.com/bwmarrin/discordgo"

	log "github.com/sirupsen/logrus"
)

type Guild struct {
	warframeChannelIds []string
	warframeRoleId     string
}

type Discord struct {
	session *discordgo.Session

	warframeGuilds []Guild
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
		session:        session,
		warframeGuilds: []Guild{},
	}
}

func (discord *Discord) Initialize(roleName string) {
	discord.initializeWarframeGuildChannels(roleName)
}

func (discord *Discord) initializeWarframeGuildChannels(roleName string) {
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

		wfGuild := Guild{
			warframeChannelIds: []string{},
		}
		for _, channel := range channels {
			log.WithField("channel", channel).Debug("channel data")

			// https://discordapp.com/developers/docs/resources/channel#channel-object-channel-types
			if channel.Type != 0 {
				continue
			}

			if isWarframeChannel(channel.Name) {
				log.WithField("channel", channel).Info("using warframe channel")
				wfGuild.warframeChannelIds = append(wfGuild.warframeChannelIds, channel.ID)
			}
		}

		if len(wfGuild.warframeChannelIds) == 0 {
			continue
		}

		if roleName != "" {
			roleId, err := discord.getMentionRoleIdByName(guild.ID, roleName)
			if err != nil {
				log.WithError(err).WithField("role", roleName).Warn("unable to @mention role")
			}

			wfGuild.warframeRoleId = roleId
		}

		discord.warframeGuilds = append(discord.warframeGuilds, wfGuild)
	}
}

func (discord *Discord) getMentionRoleIdByName(guildId string, roleName string) (string, error) {
	roles, err := discord.session.GuildRoles(guildId)
	if err != nil {
		log.WithError(err).Error("could not obtain server roles")
		return "", err
	}

	for _, role := range roles {
		log.WithField("role", role).Debug("examine role")
		if role.Name == roleName {
			log.WithField("name", roleName).WithField("id", role.ID).Info("found mention role")
			return role.ID, nil
		}
	}

	return "", errors.New("no role found with name " + roleName)
}

func (discord *Discord) Broadcast(message string) {
	for _, guild := range discord.warframeGuilds {
		var fullMessage string
		if guild.warframeRoleId != "" {
			fullMessage = "<@&" + guild.warframeRoleId + "> " + message
		} else {
			fullMessage = message
		}

		for _, channelId := range guild.warframeChannelIds {
			messageResponse, _ := discord.session.ChannelMessageSendTTS(channelId, "potato alert")

			if messageResponse != nil && messageResponse.ID != "" {
				discord.session.ChannelMessageEdit(channelId, messageResponse.ID, fullMessage)
			} else {
				discord.session.ChannelMessageSend(channelId, fullMessage)
			}
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
