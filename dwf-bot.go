package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/snakelayer/discord-warframe-orokin-alerts/discord"
	"github.com/snakelayer/discord-warframe-orokin-alerts/worldstate"

	log "github.com/sirupsen/logrus"
)

func main() {
	var token string
	var roleName string
	var isDebug bool

	flag.StringVar(&token, "token", "", "discord bot token")
	flag.StringVar(&roleName, "role", "", "role to @mention")
	flag.BoolVar(&isDebug, "debug", false, "enable debug")

	flag.Parse()

	if token == "" {
		log.Fatal("missing bot token")
	}

	if isDebug {
		log.SetLevel(log.DebugLevel)
	}

	discord := discord.New(token)
	discord.Initialize(roleName)

	ws := worldstate.New()

	go pollOrokinAlerts(ws, discord)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-ch

	discord.Close()
}

func pollOrokinAlerts(ws *worldstate.WorldState, discord *discord.Discord) {
	seenAlertIds := map[string]bool{}

	for {
		log.Info("poll worldstate")
		alerts, err := ws.GetAlerts()
		if err == nil {
			for _, alert := range alerts {
				if _, ok := seenAlertIds[alert.GetId()]; ok {
					continue
				} else {
					log.WithField("alert", alert).Info("new alert")
					seenAlertIds[alert.GetId()] = true
					discord.Broadcast(alert.PrettyPrint())
					discord.SetAlertStatus()
				}
			}

			if alerts == nil {
				log.Info("reset seen alerts")
				seenAlertIds = map[string]bool{}
				discord.ResetStatus()
			}
		}

		time.Sleep(1 * time.Minute)
	}
}
