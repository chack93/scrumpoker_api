package cleanup

import (
	"math/rand"
	"time"

	"github.com/chack93/scrumpoker_api/internal/domain/client"
	"github.com/chack93/scrumpoker_api/internal/domain/globalconfig"
	"github.com/chack93/scrumpoker_api/internal/domain/history"
	"github.com/chack93/scrumpoker_api/internal/domain/session"
	"github.com/chack93/scrumpoker_api/internal/service/msgsystem"
	"github.com/go-co-op/gocron"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func Init() error {
	msgSys := msgsystem.Get()

	msgSys.Subscribe("scrumpoker_api.cleanup", func(msg *nats.Msg) {
		logrus.Debugf("cleanup triggered")

		var sessionList []session.Session
		var err error
		if err = session.ListSession(&sessionList); err != nil {
			logrus.Errorf("cleanup, fetch session list failed, err: %v", err)
			return
		}
		expireTime := time.Now().AddDate(0, 0, -14)
		expiredSessionIDList := make([]uuid.UUID, 0)
		for _, el := range sessionList {
			if el.UpdatedAt.Before(expireTime) {
				expiredSessionIDList = append(expiredSessionIDList, el.ID)
				if err = session.DeleteSession(el.ID, &el); err != nil {
					logrus.Warnf("cleanup, delete session failed, session-id: %s, err: %v", el.ID.String(), err)
				}
			}
		}
		for _, el := range expiredSessionIDList {
			var historyList []history.History
			if err = history.ListHistoryBySessionID(el, &historyList); err != nil {
				logrus.Warnf("cleanup, list history of session failed, session-id: %s, err: %v", el.String(), err)
			}
			if historyList != nil {
				for _, hel := range historyList {
					if err = history.DeleteHistory(hel.ID, &hel); err != nil {
						logrus.Warnf("cleanup, delete history failed, history-id: %s, err: %v", hel.ID.String(), err)
					}
				}
			}
		}

		var clientList []client.Client
		if err = client.ListClient(&clientList); err != nil {
			logrus.Warnf("cleanup, list clients failed, err: %v", err)
		}
		if clientList != nil {
			for _, el := range clientList {
				if el.UpdatedAt.Before(expireTime) {
					if err = client.DeleteClient(el.ID, &el); err != nil {
						logrus.Warnf("cleanup, delete client failed, client-id: %s, err: %v", el.ID.String(), err)
					}
				}
			}
		}

		logrus.Debugf("cleanup done")
	})

	scheduler := gocron.NewScheduler(time.UTC)
	//if _, err := scheduler.Every(1).Day().At("01:00").Do(func() {
	if _, err := scheduler.Every(1).Second().Do(func() {
		// Handling multiple service instances:
		// - wait randomly of up to 30 seconds
		// - first one to run writes config & blocks others
		time.Sleep(time.Duration(rand.Int63n(30)) * time.Second)

		var lastCleanupConfig globalconfig.GlobalConfig
		if err := globalconfig.ReadGlobalConfig("last_cleanup_timestamp", &lastCleanupConfig); err != nil {
			logrus.Errorf("daily cleanup, read last cleanup time failed, err: %v", err)
			return
		}
		var lastCleanupTime time.Time
		var err error
		if lastCleanupTime, err = time.Parse(time.RFC3339, *lastCleanupConfig.Value); err != nil {
			logrus.Warnf("daily cleanup, stored last-cleanup-time invalid, continue, err: %v", err)
			lastCleanupTime = time.Now().AddDate(-1, 0, 0)
		}
		if time.Now().Sub(lastCleanupTime) > (23 * time.Hour) {
			logrus.Infof("daily cleanup, trigger cleanup")
			str := time.Now().Format(time.RFC3339)
			lastCleanupConfig.Value = &str
			globalconfig.UpdateGlobalConfig(&lastCleanupConfig)

			natsMsg := nats.NewMsg("scrumpoker_api.cleanup")
			natsMsg.Header.Add("action", "open")
			msgSys.PublishMsg(natsMsg)
		}
	}); err != nil {
		logrus.Errorf("setup scheduler failed, err: %v", err)
		return err
	}

	scheduler.StartAsync()

	return nil
}
