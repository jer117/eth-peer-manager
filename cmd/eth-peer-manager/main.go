package main

import (
	"eth-peer-manager/internal"
	"eth-peer-manager/internal/eth"
	"eth-peer-manager/internal/probber"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
	"sync"
	"time"
)

var (
	up = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "up",
		Help: "Heartbeat",
	})
	peers = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "peer_num",
		Help: "Peer counter",
	}, []string{"state"})
	gethApi         string
	probeTimeout    time.Duration
	runInterval     time.Duration
	additionalPeers string
)

func init() {
	viper.SetDefault("geth_api", "http://127.0.0.1:8545")
	viper.SetDefault("probe_timeout", 5)
	viper.SetDefault("run_interval", 60)
	viper.SetDefault("log_level", "info")
	viper.SetEnvPrefix("EPM")
	viper.AutomaticEnv()

	probeTimeout = viper.GetDuration("probe_timeout") * time.Second
	gethApi = viper.GetString("geth_api")
	runInterval = viper.GetDuration("run_interval") * time.Second
	additionalPeers = viper.GetString("additional_peers")

	// register prometheus variables
	prometheus.MustRegister(up)
	prometheus.MustRegister(peers)
	// set up to 1 by default
	up.Set(1)
	zapLogLevel := viper.GetString("log_level")
	loggerConf := zap.NewProductionConfig()
	loggerConf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	loggerConf.Level.SetLevel(utils.SetLogLevel(zapLogLevel))
	loggerConf.DisableStacktrace = true
	logger, _ := loggerConf.Build()
	zap.ReplaceGlobals(logger)
	defer logger.Sync()
}

func main() {
	totalPeerMetric := peers.With(prometheus.Labels{"state": "total"})
	unreachablePeerMetric := peers.With(prometheus.Labels{"state": "unreachable"})
	reachablePeerMetric := peers.With(prometheus.Labels{"state": "reachable"})

	var wg sync.WaitGroup
	wg.Add(2)

	//run prometheus exporter
	go func() {
		defer wg.Done()
		probber.ServeMetrics()
	}()

	for {
		if additionalPeers != "" {
			for _, ETHNode := range strings.Split(additionalPeers, ",") {
				eth.AddPeer(ETHNode, gethApi)
			}
		}

		peers := eth.GetPeers(gethApi)
		if peers != nil && len(peers) > 0 {
			unreachablePeerCount := 0
			peerCount := len(peers)

			for index, peer := range peers {
				zap.S().Debugf("index: %d, enode: %s, remote address: %s", index, peer.Enode, peer.Network.RemoteAddress)
				if !probber.Probe(peer.Network.RemoteAddress, probeTimeout) {
					eth.RemovePeer(peer.Enode, gethApi)
					unreachablePeerCount = unreachablePeerCount + 1
				}
			}

			zap.S().Infof("Peers statistics: total - %d, reachable - %d, unreachable - %d", peerCount, peerCount-unreachablePeerCount, unreachablePeerCount)
			unreachablePeerMetric.Set(float64(unreachablePeerCount))
			reachablePeerMetric.Set(float64(peerCount - unreachablePeerCount))
			totalPeerMetric.Set(float64(peerCount))
		}

		zap.S().Debugf("sleeping for %s", runInterval)
		time.Sleep(runInterval)
	}
	// Wait for the goroutines to finish.
	wg.Wait()
}
