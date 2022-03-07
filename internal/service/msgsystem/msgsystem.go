package msgsystem

import (
	"sync"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var lock = &sync.Mutex{}

type msgsystem struct {
	Conn *nats.Conn
}

var natsInstance *msgsystem

func Get() *nats.Conn {
	if natsInstance == nil {
		New()
	}
	return natsInstance.Conn
}

func New() *msgsystem {
	if natsInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if natsInstance == nil {
			natsInstance = &msgsystem{}
		}
	}
	return natsInstance
}

func (s *msgsystem) Init() error {
	natsUrl := viper.GetString("msgqueue.nats.url")
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		logrus.Errorf("connect to nats-server failed, err: %v", err)
		return err
	}
	s.Conn = conn

	logrus.Infof("connected to nats-server: %s", natsUrl)
	return nil
}
