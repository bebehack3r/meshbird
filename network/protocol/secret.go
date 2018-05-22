package protocol

import (
	"fmt"
	"io"
	"os"
  "math/rand"
	"strings"
  "time"
)

const (
	port = 3000
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
  letterIdxBits = 6                    // 6 bits to represent a letter index
  letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
  letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var (
	secMessage = RandStringBytesMaskImprSrc(64)
)

type (
	SecMessage []byte
)

func GenSecMessage() string {
	return secMessage
}

func NewSecretMessage(inp string) *Packet {
	body := Body{
		Type: TypeSec,
		Msg:  SecMessage(inp),
	}
	return &Packet{
		Head: Header{
			Length:  body.Len(),
			Version: CurrentVersion,
		},
		Data: body,
	}
}

func (s SecMessage) Len() uint16 {
	return uint16(len(s))
}

func (s SecMessage) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(s)
	return int64(n), err
}

func ReadDecodeSec(r io.Reader) (SecMessage, error) {
	logger.Debug("reading secret message...")

	secPack, errDecode := ReadAndDecode(r)
	if errDecode != nil {
		logger.Error("error on package decode, %v", errDecode)
		return nil, fmt.Errorf("error on read secret package, %v", errDecode)
	}

	if secPack.Data.Type != TypeSec {
		return nil, fmt.Errorf("non secret message received, %+v", secPack)
	}

	logger.Debug("message, %+v", secPack.Data.Msg)

	return secPack.Data.Msg.(SecMessage), nil
}

func WriteEncodeSec(w io.Writer, inp string) (err error) {
	logger.Debug("writing secret message...")
	if err = EncodeAndWrite(w, NewSecretMessage(inp)); err != nil {
		err = fmt.Errorf("error on write secret message, %v", err)
	}
	return
}

func RandStringBytesMaskImprSrc(n int) string {
  src := rand.NewSource(time.Now().UnixNano())
  b := make([]byte, n)
  for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
    if remain == 0 {
      cache, remain = src.Int63(), letterIdxMax
    }
    if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
      b[i] = letterBytes[idx]
      i--
    }
    cache >>= letterIdxBits
    remain--
  }
  return string(b)
}

func StoreSecret(fname string, msg string) (bool, error) {
	split := strings.Split(msg, ":")
	filename := fname
	formatted := fmt.Sprintf("%s:%s:%d\r\n", split[1], split[0], time.Now().Unix())
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()
	if _, err = f.WriteString(formatted); err != nil {
		return false, err
	}
	return true, nil
}
