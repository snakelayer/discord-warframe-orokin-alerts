package worldstate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	worldStateEndpoint = "http://content.warframe.com/dynamic/worldState.php"
)

var orokinItems = map[string]bool{
	"/Lotus/StoreItems/Types/Recipes/Components/OrokinCatalystBlueprint": true,
	"/Lotus/StoreItems/Types/Recipes/Components/OrokinReactorBlueprint":  true,
}

type Alert struct {
	Id struct {
		Oid string `json:"$oid"`
	} `json:"_id"`

	Activation struct {
		Date struct {
			NumberLong string `json:"$numberLong"`
		} `json:"$date"`
	} `json:"Activation"`

	Expiry struct {
		Date struct {
			NumberLong string `json:"$numberLong"`
		} `json:"$date"`
	} `json:"Expiry"`

	MissionInfo struct {
		Difficulty     float64 `json:"difficulty"`
		EnemySpec      string  `json:"enemySpec"`
		ExtraEnemySpec string  `json:"extraEnemySpec"`
		Faction        string  `json:"faction"`
		LevelOverride  string  `json:"levelOverride"`
		Location       string  `json:"location"`
		MaxEnemyLevel  int     `json:"maxEnemyLevel"`
		MaxWaveNum     int     `json:"maxWaveNum"`
		MinEnemyLevel  int     `json:"minEnemyLevel"`
		MissionReward  struct {
			Credits int      `json:"credits"`
			Items   []string `json:"items"`
		} `json:"missionReward"`
		MissionType string `json:"missionType"`
		Seed        int    `json:"seed"`
	} `json:"MissionInfo"`
}

func (alert *Alert) String() string {
	return fmt.Sprintf("{Id:%v}", alert.GetId())
}

func (alert *Alert) GetId() string {
	return alert.Id.Oid
}

func (alert *Alert) PrettyPrint() string {
	var buffer bytes.Buffer

	buffer.WriteString("**" + WorldStateItems[alert.MissionInfo.MissionReward.Items[0]] + "**")
	//buffer.WriteString("**" + alert.MissionInfo.MissionReward.Items[0] + "**")
	buffer.WriteString(": ")
	buffer.WriteString(WorldStateLocations[alert.MissionInfo.Location])
	buffer.WriteString(" | ")
	buffer.WriteString(WorldStateMissionTypes[alert.MissionInfo.MissionType])
	buffer.WriteString(" | ")
	buffer.WriteString(WorldStateFaction[alert.MissionInfo.Faction])
	buffer.WriteString(" ")
	buffer.WriteString(strconv.Itoa(alert.MissionInfo.MinEnemyLevel))
	buffer.WriteString("-")
	buffer.WriteString(strconv.Itoa(alert.MissionInfo.MaxEnemyLevel))
	buffer.WriteString(" *(" + getMinutesUntil(alert.Expiry.Date.NumberLong).String() + ")*")

	return buffer.String()
}

func getMinutesUntil(expire string) time.Duration {
	expireMillis, err := strconv.ParseInt(expire, 10, 64)
	if err != nil {
		log.WithError(err).WithField("expire", expire).Error("could not convert time")
		return -1
	}

	expireDate := time.Unix(expireMillis/1000, 0)

	return expireDate.Sub(time.Now()).Round(time.Minute)
}

type WorldState struct {
	Alerts []*Alert `json:"Alerts"`
}

func New() *WorldState {
	return new(WorldState)
}

func (ws *WorldState) refresh() error {
	resp, err := http.Get(worldStateEndpoint)
	if err != nil {
		log.WithError(err).Error("request to worldstate API failedbv")
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(ws)
	if err != nil {
		log.WithError(err).Error("error decoding worldstate data to json")
		return err
	}

	return nil
}

func (ws *WorldState) getOrokinAlerts() []*Alert {
	orokinAlerts := []*Alert{}
	for _, alert := range ws.Alerts {
		missionReward := alert.MissionInfo.MissionReward

		for _, item := range missionReward.Items {
			if _, ok := orokinItems[item]; ok {
				log.WithField("alert", alert).Info("found orokin alert")
				orokinAlerts = append(orokinAlerts, alert)
				break
			}
		}
	}

	return orokinAlerts
}

func (ws *WorldState) GetAlerts() ([]*Alert, error) {
	err := ws.refresh()
	if err != nil {
		log.WithError(err).Error("could not refresh worldstate")
		return nil, err
	}

	return ws.getOrokinAlerts(), nil
}
