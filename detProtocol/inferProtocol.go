package detProtocol

import "regexp"

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

const (
	Local = iota
	HTTP  = iota
)
