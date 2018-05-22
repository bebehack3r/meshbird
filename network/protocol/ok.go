package protocol

import (
	"fmt"
	"io"

	"github.com/rmnoff/meshbird/secure"
)

var (
	onMessage = []byte{'O', 'K'}
)

type (
	OkMessage []byte
)

func NewOkMessage(session []byte) *Packet {
	self, _ := secure.GetSelf(3003)
	first := append(onMessage, splitter...)
	second := append(first, []byte(self)...)
	tmp := append(second, splitter...)
	body := Body{
		Type: TypeOk,
		Msg:  OkMessage(append(tmp, session...)),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (o OkMessage) Len() uint16 {
	return uint16(len(o))
}

func (o OkMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(o)
	return int64(n), err
}

func (o OkMessage) SessionKey() []byte {
	before := len(onMessage) + len(splitter) * 2 + len([]byte("0x0000000000000000000000000000000000000000"))
	return o[before:]
}

func (o OkMessage) Address() []byte {
	before := len(onMessage) + len(splitter)
	address := []byte("0x0000000000000000000000000000000000000000")
	completeLength := before + len(address)
	after := before + completeLength
	return o[before:after]
}

func ReadDecodeOk(r io.Reader) (OkMessage, error) {
	logger.Debug("reading ok message...")

	okPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Error("error on package decode, %v", errDecode)
		return nil, fmt.Errorf("error on read ok package, %v", errDecode)
	}

	if okPack.Data.Type != TypeOk {
		return nil, fmt.Errorf("non ok message received, %+v", okPack)
	}

	logger.Debug("message, %v", okPack.Data.Msg)
	return okPack.Data.Msg.(OkMessage), nil
}

func WriteEncodeOk(w io.Writer, session []byte) (err error) {
	logger.Debug("writing ok message...")
	if err = EncodeAndWrite(w, NewOkMessage(session)); err != nil {
		err = fmt.Errorf("error on write ok message, %v", err)
	}
	return
}
