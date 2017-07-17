package rpc

import (
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/skycoin/net/skycoin-messenger/factory"
	"github.com/skycoin/net/skycoin-messenger/msg"
	"github.com/skycoin/skycoin/src/cipher"
)

var DefaultClient = &Client{Push: make(chan interface{}, 8)}

type Client struct {
	sync.RWMutex
	Connection *factory.Connection

	Push   chan interface{}
	Logger *log.Entry
}

func (c *Client) GetConnection() *factory.Connection {
	c.RLock()
	defer c.RUnlock()
	return c.Connection
}

func (c *Client) SetConnection(connection *factory.Connection) {
	c.Lock()
	if c.Connection != nil {
		c.Connection.Close()
	}
	c.Connection = connection
	c.Unlock()
}

func (c *Client) PushLoop(conn *factory.Connection) {
	defer func() {
		if err := recover(); err != nil {
			c.Logger.Errorf("PushLoop recovered err %v", err)
		}
	}()
	push := &msg.PushMsg{}
	for {
		select {
		case m, ok := <-conn.GetChanIn():
			if !ok || len(m) < factory.MSG_HEADER_END {
				return
			}
			op := m[factory.MSG_OP_BEGIN]
			switch op {
			case factory.OP_SEND:
				if len(m) < factory.MSG_META_END {
					continue
				}
				key := cipher.NewPubKey(m[factory.MSG_PUBLIC_KEY_BEGIN:factory.MSG_PUBLIC_KEY_END])
				push.From = key.Hex()
				push.Msg = string(m[factory.MSG_META_END:])
				c.Push <- push
			case factory.OP_REG:
				if len(m) < factory.MSG_PUBLIC_KEY_END {
					continue
				}
				key := cipher.NewPubKey(m[factory.MSG_PUBLIC_KEY_BEGIN:factory.MSG_PUBLIC_KEY_END])
				c.Connection.SetKey(key)
				c.Push <- &msg.Reg{PublicKey:key.Hex()}
			}
		}
	}
}
