package socketmsg

import (
	"encoding/json"

	"github.com/chack93/scrumpoker_api/internal/domain/client"
	"github.com/chack93/scrumpoker_api/internal/domain/history"
	"github.com/chack93/scrumpoker_api/internal/domain/session"
)

type SocketMsgHead struct {
	Action   string `json:"action"`
	ClientID string `json:"clientId"`
	GroupID  string `json:"groupId"`
}
type SocketMsgBodyUpdate struct {
	Session     *session.Session   `json:"session"`
	Client      *client.Client     `json:"client"`
	ClientList  *[]client.Client   `json:"clientList"`
	HistoryList *[]history.History `json:"historyList"`
}
type SocketMsg struct {
	Head SocketMsgHead   `json:"head"`
	Body json.RawMessage `json:"body"`
}
