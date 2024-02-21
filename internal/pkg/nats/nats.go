package nats

import (
	"time"

	"github.com/nats-io/nats.go"
)

const (
	maxReconnects   = 60
	reconnectWaitMs = 2000
)

func NewNatsConn(natsServer string) (*nats.Conn, error) {
	return nats.Connect(
		natsServer,
		nats.ReconnectWait(time.Millisecond*time.Duration(reconnectWaitMs)),
		nats.MaxReconnects(maxReconnects),
	)
}
