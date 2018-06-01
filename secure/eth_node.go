package secure

import (
  "bufio"
  "encoding/json"
  "encoding/hex"
  "fmt"
  "io/ioutil"
  "net/http"
  "os"
  "os/exec"
  "strings"
  "bytes"
)

type HTTPResponse struct {
  Success bool `json:"success"`
  Response string `json:"result"`
}

func Encode(port int, cipher string, message string) (result string, err error) {
  err = nil
  url := fmt.Sprintf("http://localhost:%d/enc", port)

  var jsonStr = []byte(fmt.Sprintf(`{"message":"%s", "goCipher":"%s"}`, message, cipher))
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Close = true
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    fmt.Printf("%s", err)
    return "", err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  node := HTTPResponse{}
  json.Unmarshal(body, &node)
  self := node.Response
  return self, nil
}

func CheckInbox(port int, s string) (result string, err error) {
  err = nil
  url := fmt.Sprintf("http://localhost:%d/inbox", port)

  var jsonStr = []byte(fmt.Sprintf(`{"s": "%s"}`, s))
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Close = true
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  node := HTTPResponse{}
  json.Unmarshal(body, &node)
  msgs := node.Response
  return msgs, nil
}

func SendMessage(port int, address string, message string) (result string, err error) {
  err = nil
  url := fmt.Sprintf("http://localhost:%d/", port)

  var jsonStr = []byte(fmt.Sprintf(`{"receiver":"%s", "message":"%s"}`, address, message))
  req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
  req.Close = true
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, _ := ioutil.ReadAll(resp.Body)
  node := HTTPResponse{}
  json.Unmarshal(body, &node)
  self := node.Response
  return self, nil
}

func GetSelf(port int) (self string, err error) {
  err = nil
  resp, err := http.Get(fmt.Sprintf("http://localhost:%d/self", port))
  if err != nil {
    return "", err
  }
  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return "", err
  }
  node := HTTPResponse{}
  json.Unmarshal(body, &node)
  self = node.Response
  return
}

func GetHexedSelf(port int) (hexedPart string, err error) {
  self, err := GetSelf(port)
  hexedPart = hex.EncodeToString([]byte(self))
  return
}

func RunNode(port int, address string, privateKey string) *exec.Cmd {
  cmd := exec.Command("node", "index.js")
  env := os.Environ()
  env = append(env, fmt.Sprintf("PORT=%d", port))
  env = append(env, fmt.Sprintf("ADDRESS=%s", address))
  env = append(env, fmt.Sprintf("PRIVATE_KEY=%s", privateKey))
  env = append(env, fmt.Sprintf("INFURA_CONNECTION=https://mainnet.infura.io/4mNGIiR3CMs8xz0uLvJE"))
  cmd.Env = env
  go cmd.Start()
  return cmd
}

func NodeAuthorize() (string, string) {
  fmt.Printf("Please, authorize using your ETH wallet address and private key.\nYour address: ")
  reader := bufio.NewReader(os.Stdin)
  walletAddress, _ := reader.ReadString('\n')
  walletAddress = strings.Replace(walletAddress, "\n", "", -1)
  fmt.Printf("Your private key: ")
  reader = bufio.NewReader(os.Stdin)
  walletKey, _ := reader.ReadString('\n')
  walletKey = strings.Replace(walletKey, "\n", "", -1)
  return walletAddress, walletKey
}
