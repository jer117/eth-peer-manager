package eth

import (
	"encoding/json"
	"eth-peer-manager/internal"
	"go.uber.org/zap"
)

func GetBlock(gethApi string) string {
	resp, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"method":"eth_blockNumber","params":[],"id":1,"jsonrpc":"2.0"}`).
		Post(gethApi)
	if err != nil {
		zap.S().Errorf("failed to get current block number: %s", err)
	}
	if resp.Body() != nil {
		body := make(map[string]interface{})
		err := json.Unmarshal(resp.Body(), &body)
		if err != nil {
			zap.S().Errorf("failed to unmarshal blockNumber call response: %s", err)
			return ""
		}
		return body["result"].(string)
	}
	return ""
}
