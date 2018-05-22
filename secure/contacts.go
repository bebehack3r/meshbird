package secure

import (
	"bufio"
	"os"
	"strings"
)

func FindCipher(address string) []byte {
  f, err := os.Open("db")
  if err != nil {
    return nil
  }
  defer f.Close()

  scanner := bufio.NewScanner(f)

  for scanner.Scan() {
    if strings.Contains(scanner.Text(), address) {
      return []byte(strings.Split(scanner.Text(), ":")[0])
    }
  }

  return nil
}
