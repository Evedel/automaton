package processor

import (
	"actioneer/internal/state"
	"encoding/json"
	"fmt"
	"log/slog"
)

var alertNameLabel = "alertname"

type Alert struct {
	Status string
	Labels map[string]string
}

type Notification struct {
	Alerts []Alert
}

func ReadIncommingNotification(bytes []byte) (Notification, error) {
	slog.Debug("incomming bytes: " + fmt.Sprintf("%+v", string(bytes)))
	var notification Notification
	err := json.Unmarshal(bytes, &notification)
	if err != nil {
		slog.Error("cannot unmarshal incomming bytes: " + fmt.Sprintf("%+v", string(bytes)))
		slog.Error(err.Error())
	}
	return notification, err
}

func CheckActionNeeded(state state.State, notification Notification) bool {
	slog.Debug("incomming notification: " + fmt.Sprint(notification))

	if len(notification.Alerts) == 0 {
		slog.Error("no alerts in notification: " + fmt.Sprint(notification))
		return false
	}

	for _, alert := range notification.Alerts {
		if alertName, ok := alert.Labels[alertNameLabel]; ok {
			_, found := state.GetActionByAlertName(alertName)
			if found && (alert.Status == "firing") {
				return true
			}
		}
	}
	return false
}
