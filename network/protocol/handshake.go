package protocol

import (
	"bytes"
	"fmt"
	"github.com/rmnoff/meshbird/secure"
	"io"
)

var (
	magicKey = []byte{'C', 'P', 'H', 'R'}
	splitter = []byte{':'}
)

type (
	HandshakeMessage []byte
)

func IsMagicValid(data []byte) bool {
	logger.Debug("decoding magic... expected %v, got %v", magicKey, data)
	return bytes.HasPrefix(data, magicKey)
}

func NewHandshakePacket(sessionKey []byte, networkSecret *secure.NetworkSecret, address []byte) *Packet {
	first := append(magicKey, splitter...)
	second := append(first, address...)
	tmp := append(second, splitter...)
	sessionKey = append(tmp, sessionKey...)
	data := networkSecret.Encode(sessionKey)

	body := Body{
		Type: TypeHandshake,
		Msg:  HandshakeMessage(data),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (m HandshakeMessage) Len() uint16 {
	return uint16(len(m))
}

func (m HandshakeMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(m)
	return int64(n), err
}

func (m HandshakeMessage) Bytes() []byte {
	return []byte(m)
}

func (m HandshakeMessage) SessionKey() []byte {
	length := len(magicKey) + len(splitter) * 2 + len([]byte("0x0000000000000000000000000000000000000000"))
	return m[length:]
}

func (m HandshakeMessage) Address() []byte {
	before := len(magicKey) + len(splitter)
	after := len(splitter) + len([]byte("0x0000000000000000000000000000000000000000"))
	return m[before:after]
}

func ReadDecodeHandshake(r io.Reader) (HandshakeMessage, error) {
	logger.Debug("reading handshare message...")

	handshakePack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Error("error on package decode, %v", errDecode)
		return nil, fmt.Errorf("error on read handshare package, %v", errDecode)
	}

	if handshakePack.Data.Type != TypeHandshake {
		return nil, fmt.Errorf("non handshare message received, %+v", handshakePack)
	}

	logger.Debug("message, %v", handshakePack.Data.Msg)
	return handshakePack.Data.Msg.(HandshakeMessage), nil
}

func WriteEncodeHandshake(w io.Writer, sessionKey []byte, networkSecret *secure.NetworkSecret, address []byte) (err error) {
	logger.Debug("writing handshare message...")
	if err = EncodeAndWrite(w, NewHandshakePacket(sessionKey, networkSecret, address)); err != nil {
		err = fmt.Errorf("error on write handshare message, %v", err)
	}
	return
}
