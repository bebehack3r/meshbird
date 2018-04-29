package common

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rmnoff/meshbird/log"
	"github.com/rmnoff/meshbird/network/protocol"
	"github.com/rmnoff/meshbird/secure"
)

var (
	rnLoggerFormat = "remote %s"
)

type RemoteNode struct {
	Node
	conn             net.Conn
	sessionKey       []byte
	privateIP        net.IP
	publicAddress    string
	logger           log.Logger
	lastHeartbeat    time.Time
	currentSecretKey string
}

func NewRemoteNode(conn net.Conn, sessionKey []byte, privateIP net.IP) *RemoteNode {
	return &RemoteNode{
		conn:          conn,
		sessionKey:    sessionKey,
		privateIP:     privateIP,
		publicAddress: conn.RemoteAddr().String(),
		logger:        log.L(fmt.Sprintf(rnLoggerFormat, privateIP.String())),
		lastHeartbeat: time.Now(),
	}
}

func (rn *RemoteNode) SendToInterface(payload []byte) error {
	return protocol.WriteEncodeTransfer(rn.conn, payload)
}

func (rn *RemoteNode) SendPack(pack *protocol.Packet) (err error) {
	if err = protocol.EncodeAndWrite(rn.conn, pack); err != nil {
		err = fmt.Errorf("error on write transfer message, %v", err)
	}
	return
}

func (rn *RemoteNode) Close() {
	defer rn.conn.Close()
	rn.logger.Debug("closing...")
}

func (rn *RemoteNode) listen(ln *LocalNode) {
	defer rn.logger.Debug("listener stopped...")
	defer func() {
		ln.NetTable().RemoveRemoteNode(rn.privateIP)
	}()

	iface, ok := ln.Service("iface").(*InterfaceService)
	if !ok {
		rn.logger.Error("interface service not found")
		return
	}

	rn.logger.Debug("listening...")

	for {
		pack, err := protocol.Decode(rn.conn)
		if err != nil {
			rn.logger.Error("decode error, %v", err)
			if err == io.EOF {
				break
			}
			continue
		}
		rn.logger.Debug("received, %+v", pack)

		switch pack.Data.Type {
		case protocol.TypeTransfer:
			rn.logger.Debug("Writing to interface...")
			payloadEncrypted := pack.Data.Msg.(protocol.TransferMessage).Bytes()
			payload, errDec := secure.DecryptIV(payloadEncrypted, ln.State().Secret.Key)
			if errDec != nil {
				rn.logger.Error("error on decrypt, %v", err)
				break
			}
			srcAddr := net.IP(payload[12:16])
			dstAddr := net.IP(payload[16:20])
			rn.logger.Info("received packet from %s to %s", srcAddr.String(), dstAddr.String())
			if err = iface.WritePacket(payload); err != nil {
				rn.logger.Error("write packet err: %s", err)
			}
		case protocol.TypeHeartbeat:
			rn.logger.Debug("heartbeat received, %v", pack.Data.Msg)
			rn.lastHeartbeat = time.Now()
		case protocol.TypeSec:
			rn.logger.Info("Secret received, here it is: %s", pack.Data.Msg)
			rn.StoreSecret(pack.Data.Msg)
		}
	}
}

func (rn *RemoteNode) StoreSecret(msg protocol.Message) {
	fmt.Printf("This secret: %s\nis now stored for:\n%v\n\n\n", msg, string(rn.sessionKey))
}
