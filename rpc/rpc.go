package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

type BaseMessage struct {
	Method string `json:"method"`
}

func GetBaseMessage(msg []byte) ([]byte, []byte, int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if !found {
		return nil, nil, 0, nil
	}

	contentLengthBytes := header[len("Content-Length: "):]
	contentLength, err := strconv.Atoi(string(contentLengthBytes))
	if err != nil {
		return nil, nil, 0, err
	}

	return header, content, contentLength, nil
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, contentLength, err := GetBaseMessage(msg)
	_ = header // Not required for this function aas only content is needed
	if err != nil {
		return "n/a", nil, err
	}

	var baseMessage BaseMessage
	if err := json.Unmarshal(content[:contentLength], &baseMessage); err != nil {
		return "n/a", nil, nil
	}

	return baseMessage.Method, content[:contentLength], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, contentLength, err := GetBaseMessage(data)
	if err != nil {
		return 0, nil, err
	}

	if len(content) < contentLength {
		return 0, nil, nil
	}

	//4 comes from '\r''\n''\r''\n'
	totalLength := len(header) + 4 + contentLength

	// println("Message length from my way: %d", len(data))

	return totalLength, data[:totalLength], nil
}
