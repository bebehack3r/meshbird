package secure

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"encoding/base64"
	"io/ioutil"
)

func Compress(input string) string {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(input)); err != nil {
		panic(err)
	}
	if err := gz.Flush(); err != nil {
		panic(err)
	}
	if err := gz.Close(); err != nil {
		panic(err)
	}
	str := base64.StdEncoding.EncodeToString(b.Bytes())
	fmt.Println(str)
	data, _ := base64.StdEncoding.DecodeString(str)
	fmt.Println(data)
	rdata := bytes.NewReader(data)
	r,_ := gzip.NewReader(rdata)
	s, _ := ioutil.ReadAll(r)
	return string(s)
}
