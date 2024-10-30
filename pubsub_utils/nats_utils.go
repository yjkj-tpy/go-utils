package pubsub_utils

import (
	"github.com/nats-io/nats.go"
	"log"
	"sync"
)

type NatsCfg struct {
	NatsUrl string `json:"natsUrl"`
}

var (
	natsOnce     sync.Once
	natsInstance *nats.Conn
)

// GetNatsServer 获取单例的 NatsServer 连接
func GetNatsServer(cfg *NatsCfg) *nats.Conn {
	natsOnce.Do(func() {
		nc, err := nats.Connect(cfg.NatsUrl)
		if err != nil {
			log.Fatalf("Error connecting to NATS: %v", err)
		}
		natsInstance = nc
	})

	return natsInstance
}
