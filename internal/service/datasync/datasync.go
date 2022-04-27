package datasync

import (
	"encoding/json"
	"fmt"

	"github.com/chack93/scrumpoker_api/internal/domain/client"
	"github.com/chack93/scrumpoker_api/internal/domain/history"
	"github.com/chack93/scrumpoker_api/internal/domain/session"
	"github.com/chack93/scrumpoker_api/internal/domain/socketmsg"
	"github.com/chack93/scrumpoker_api/internal/service/msgsystem"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func Init() error {
	msgSys := msgsystem.Get()

	msgSys.Subscribe("scrumpoker_api.client-request", func(msg *nats.Msg) {
		action := msg.Header.Get("action")
		clientID, err := uuid.Parse(msg.Header.Get("clientID"))
		if err != nil {
			logrus.Errorf("invalid cliend-UUID, cID: %s, err: %v", clientID.String(), err)
			return
		}
		groupID, err := uuid.Parse(msg.Header.Get("groupID"))
		if err != nil {
			logrus.Errorf("invalid group-UUID, gID: %s, err: %v", groupID.String(), err)
			return
		}

		logrus.Debugf("action: %s, cID: %s, gID: %s", action, clientID.String(), groupID)

		switch action {
		case "open":
			var cl client.Client
			if err := client.ReadClient(clientID, &cl); err != nil {
				logrus.Errorf("read client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
				return
			}
			t := true
			gid := groupID.String()
			cl.Connected = &t
			cl.SessionId = &gid
			if err := client.UpdateClient(clientID, &cl); err != nil {
				logrus.Errorf("update client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
				return
			}
			UpdateClientsOfGroup(groupID)
		case "close":
			var cl client.Client
			if err := client.ReadClient(clientID, &cl); err != nil {
				logrus.Errorf("read client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
				return
			}
			t := false
			cl.Connected = &t
			if err := client.UpdateClient(clientID, &cl); err != nil {
				logrus.Errorf("update client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
				return
			}
			UpdateClientsOfGroup(groupID)
		case "update":
			handleUpdateRequest(msg)
		default:
			logrus.Errorf("unknown action in request, action: %s, cID: %s", action, clientID.String())
			return
		}
	})
	return nil
}

func UpdateClientsOfGroup(groupID uuid.UUID) (err error) {
	var se session.Session
	if err = session.ReadSession(groupID, &se); err != nil {
		logrus.Errorf("read session failed, gID: %s, err: %v", groupID.String(), err)
		return
	}
	var clList []client.Client
	if err = client.ListClientOfSession(groupID, &clList); err != nil {
		logrus.Errorf("read client list failed, gID: %s, err: %v", groupID.String(), err)
		return
	}
	var hiList []history.History
	if err = history.ListHistoryBySessionID(groupID, &hiList); err != nil {
		logrus.Errorf("read client list failed, gID: %s, err: %v", groupID.String(), err)
		return
	}
	msgBody := socketmsg.SocketMsgBodyUpdate{
		Session:     &se,
		ClientList:  &clList,
		HistoryList: &hiList,
	}

	for _, el := range clList {
		msgSys := msgsystem.Get()
		natsMsg := nats.NewMsg(fmt.Sprintf("scrumpoker_api.client-response.%s", el.ID.String()))
		natsMsg.Header.Add("clientID", el.ID.String())
		natsMsg.Header.Add("groupID", groupID.String())
		natsMsg.Header.Add("action", "update")

		msgBody.Client = &el
		bodyJson, err := json.Marshal(msgBody)
		if err != nil {
			logrus.Errorf("marshal update body failed, gID: %s, err: %v", groupID.String(), err)
			continue
		}
		natsMsg.Data = bodyJson
		msgSys.PublishMsg(natsMsg)
	}
	return nil
}

func handleUpdateRequest(msg *nats.Msg) {
	action := msg.Header.Get("action")
	clientID, err := uuid.Parse(msg.Header.Get("clientID"))
	if err != nil {
		logrus.Errorf("invalid cliend-UUID, cID: %s, err: %v", clientID.String(), err)
		return
	}
	groupID, err := uuid.Parse(msg.Header.Get("groupID"))
	if err != nil {
		logrus.Errorf("invalid group-UUID, gID: %s, err: %v", groupID.String(), err)
		return
	}
	var updateRequest socketmsg.SocketMsgBodyUpdate
	if err := json.Unmarshal(msg.Data, &updateRequest); err != nil {
		logrus.Errorf("unmarshal client update failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
		return
	}
	var se session.Session
	if err := session.ReadSession(groupID, &se); err != nil {
		logrus.Errorf("read session failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
		return
	}
	var cl client.Client
	if err := client.ReadClient(clientID, &cl); err != nil {
		logrus.Errorf("read client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
		return
	}

	if updateRequest.Client != nil {
		connectedTrue := true
		cl.Connected = &connectedTrue
		cl.Estimation = updateRequest.Client.Estimation
		cl.Name = updateRequest.Client.Name
		cl.SessionId = updateRequest.Client.SessionId
		if cl.SessionId == nil || len(*cl.SessionId) < 1 {
			grpId := groupID.String()
			cl.SessionId = &grpId
		}
		cl.Viewer = updateRequest.Client.Viewer
	}

	if se.OwnerClientId != nil && *se.OwnerClientId == cl.ID.String() && updateRequest.Session != nil {
		oldGameStatus := ""
		if se.GameStatus != nil {
			oldGameStatus = *se.GameStatus
		}
		se.CardSelectionList = updateRequest.Session.CardSelectionList
		se.Description = updateRequest.Session.Description
		//se.OwnerClientId = updateRequest.Session.OwnerClientId
		se.GameStatus = updateRequest.Session.GameStatus
		if err := session.UpdateSession(clientID, &se); err != nil {
			logrus.Errorf("update client failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
			return
		}

		if oldGameStatus != "new" && updateRequest.Session.GameStatus != nil && *updateRequest.Session.GameStatus == "new" {
			var clientList []client.Client
			if err := client.ListClientOfSession(groupID, &clientList); err != nil {
				logrus.Errorf("update failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
			}

			for _, el := range clientList {
				emptyEstimation := ""
				el.Estimation = &emptyEstimation
				cl.Estimation = &emptyEstimation

				if err := client.UpdateClient(el.ID, &el); err != nil {
					logrus.Errorf("new game / reset client estimation failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
				}
			}
		}

		if oldGameStatus != "reveal" && updateRequest.Session.GameStatus != nil && *updateRequest.Session.GameStatus == "reveal" {
			gameUUID := uuid.New().String()

			var clientList []client.Client
			if err := client.ListClientOfSession(groupID, &clientList); err != nil {
				logrus.Errorf("update failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
			}

			for _, el := range clientList {
				cid := el.ID.String()
				gid := groupID.String()
				if err := history.CreateHistory(&history.History{
					HistoryNew: history.HistoryNew{
						ClientId:   &cid,
						ClientName: el.Name,
						Estimation: el.Estimation,
						SessionId:  &gid,
					},
					GameId: &gameUUID,
				}); err != nil {
					logrus.Errorf("add history item failed, action: %s, cID: %s, gID: %s, err: %v", action, clientID.String(), groupID.String(), err)
				}
			}
		}
	}

	if err := client.UpdateClient(clientID, &cl); err != nil {
		logrus.Errorf("update client failed, action: %s, cID: %s, err: %v", action, clientID.String(), err)
		return
	}
	UpdateClientsOfGroup(groupID)
}
