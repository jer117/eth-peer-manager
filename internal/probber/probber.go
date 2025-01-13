package probber

import (
	"go.uber.org/zap"
	"net"
	"time"
)

func Probe(host string, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", host, timeout)
	if err != nil {
		zap.S().Infof("%s: %s", host, err)
		return false
	}
	if conn != nil {
		defer conn.Close()
		zap.S().Debugf("%s: connection established", host)
		return true
	}
	return true
}
