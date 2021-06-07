package detProtocol

import "regexp"

const (
	Local = iota
	HTTP  = iota
)

func InferProtocolFromFilename(filename string) int {
	matched, err := regexp.Match(`^https?:\/\/`, []byte(filename))
	if err != nil {
		panic(err)
	}
	if matched {
		return HTTP
	} else {
		return Local
	}
}
