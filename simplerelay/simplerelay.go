package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/timesking/livego/logs"
	"github.com/timesking/livego/protocol/rtmp/rtmprelay"
	"go.uber.org/zap"
)

var (
	version = "master"

	rtmpPullAddr     = flag.String("rtmp-source-addr", "rtmp://127.0.0.1:1935", "RTMP server source/pull from")
	rtmpPushAddr     = flag.String("rtmp-dest-addr", "rtmp://xxx.xxx.xxx.xxx", "RTMP server dest/push to")
	rtmpPushIPRelace = flag.String("rtmp-dest-ip", "", "RTMP server ip assigned manually")
)

func init() {

	flag.Parse()
}

func main() {
	logger, _ = zap.NewProduction()
	logger = logger.Named("main")
	defer logger.Sync()
	zaplog = logger.Sugar()
	zaplog.Info("simple relay rtmp:", version)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	localurl := *rtmpPullAddr
	remoteurl := *rtmpPushAddr

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		cancel()
	}()
	logs.SetLoger(NewRelayLoger())
	logs.SetZapLogerBuilder(func(options ...zap.Option) *zap.Logger {
		return logger.WithOptions(options...)
	})
	pushRtmprelay := rtmprelay.NewRtmpRelay(&localurl, &remoteurl)

	if rtmpPushIPRelace != nil && len(*rtmpPushIPRelace) != 0 {
		pushRtmprelay.DNSLookup = func(host string) (string, error) {
			return *rtmpPushIPRelace, nil
		}
	}

	zaplog.Infof("rtmprelay start push %s from %s", remoteurl, localurl)
	err := pushRtmprelay.StartWait(ctx)

	if err != nil {
		zaplog.Errorf("push error=%v", err)
		logger.Sync()
		os.Exit(-1)
	}
}
