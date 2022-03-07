package datasync

import (
	"fmt"

	"github.com/chack93/scrumpoker_api/internal/service/msgsystem"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

/*
type WebsocketSessionInfo struct {
	sessionID             string
	connectedClientIDList []string
}

// TODO: put this info in db
var sessionMap map[string]*WebsocketSessionInfo

func Get(sessionID string) *WebsocketSessionInfo {
	if sessionMap[sessionID] == nil {
		sessionMap[sessionID] = &WebsocketSessionInfo{
			sessionID:             sessionID,
			connectedClientIDList: []string{},
		}
	}
	return sessionMap[sessionID]
}
*/

func Init() error {
	msgSys := msgsystem.Get()

	msgSys.Subscribe("scrumpoker_api.client-request", func(msg *nats.Msg) {
		action := msg.Header.Get("action")
		clientID := msg.Header.Get("clientID")
		groupID := msg.Header.Get("groupID")
		logrus.Debugf("action: %s, cID: %s, gID: %s", action, clientID, groupID)

		natsMsg := nats.NewMsg(fmt.Sprintf("scrumpoker_api.client-response.%s", clientID))
		natsMsg.Header.Add("clientID", clientID)
		natsMsg.Header.Add("groupID", groupID)
		natsMsg.Header.Add("action", action)
		natsMsg.Data = msg.Data
		msgSys.PublishMsg(natsMsg)
	})
	return nil
}
