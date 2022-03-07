package fabric

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
)

type Message struct {
	Log string `json:"log"`
}

func Create(info *InitInfo, key, value string) error {
	messageOb := Message{Log: value}
	message, err := json.Marshal(messageOb)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err.Error())
	}
	req := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn:         "create",
		Args:        [][]byte{[]byte(key), message},
	}
	_, err = info.ChannelClient.Execute(req, channel.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return fmt.Errorf("failed to invoke: %v", err.Error())
	}
	return nil
}

func Query(info *InitInfo, key string) (string, error) {
	request := channel.Request{
		ChaincodeID: info.ChaincodeID,
		Fcn:         "query",
		Args:        [][]byte{[]byte(key)},
	}
	response, err := info.ChannelClient.Query(request, channel.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return "", fmt.Errorf("failed to query: %v", err.Error())
	}
	return string(response.Payload), nil
}
