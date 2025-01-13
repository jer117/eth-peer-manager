package eth

import (
	"encoding/json"
	"eth-peer-manager/internal"
	"fmt"
	"go.uber.org/zap"
)

type Peers struct {
	Jsonrpc string `json:"jsonrpc"`
	Id      int    `json:"id"`
	Result  []Peer `json:"result"`
}

type Peer struct {
	Enode   string   `json:"enode"`
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Caps    []string `json:"caps"`
	Network struct {
		LocalAddress  string `json:"localAddress"`
		RemoteAddress string `json:"remoteAddress"`
		Inbound       bool   `json:"inbound"`
		Trusted       bool   `json:"trusted"`
		Static        bool   `json:"static"`
	} `json:"network"`
	Protocols struct {
		Eth struct {
			Version int `json:"version"`
		} `json:"eth"`
	} `json:"protocols"`
	Enr string `json:"enr,omitempty"`
}

func GetPeers(gethApi string) []Peer {
	list := &Peers{}
	resp, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(`{"method":"admin_peers","params":[],"id":1,"jsonrpc":"2.0"}`).
		SetResult(list).
		Post(gethApi)
	if err != nil {
		zap.S().Errorf("failed to get peer list: %s", err)
	}
	if resp.Body() != nil {
		return list.Result
	}
	return nil
}

func RemovePeer(peer string, gethApi string) {
	resp, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"method":"admin_removePeer","params":["%s"],"id":1,"jsonrpc":"2.0"}`, peer)).
		Post(gethApi)
	if err != nil {
		zap.S().Errorf("failed to remove peer: %s", err)
	}
	if resp != nil {
		zap.S().Debugf("removePeer call response: %s", resp)
	}
}

func AddPeer(peer string, gethApi string) {
	resp, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"method":"admin_addPeer","params":["%s"],"id":1,"jsonrpc":"2.0"}`, peer)).
		Post(gethApi)
	if err != nil {
		zap.S().Errorf("failed to add peer: %s", err)
	}
	if resp != nil {
		body := make(map[string]interface{})
		err := json.Unmarshal(resp.Body(), &body)
		if err != nil {
			zap.S().Errorf("failed to unmarshal addPeer call response: %s", err)
		}
		zap.S().Debugf("addPeer call response: %s", resp)
	}
}
