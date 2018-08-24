package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/timesking/livego/protocol/rtmp/rtmprelay"
	"go.uber.org/zap"
)

var (
	version      = "master"
	logger       *zap.Logger
	zaplog       *zap.SugaredLogger
	done         chan bool
	rtmpPullAddr = flag.String("rtmp-source-addr", "rtmp://127.0.0.1:1935", "RTMP server source/pull from")
	rtmpPushAddr = flag.String("rtmp-dest-addr", "rtmp://xxx.xxx.xxx.xxx", "RTMP server dest/push to")
)

func init() {
	done = make(chan bool, 1)
	flag.Parse()
}

func main() {
	logger, _ = zap.NewProduction()
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

	pushRtmprelay := rtmprelay.NewRtmpRelay(&localurl, &remoteurl)
	pushRtmprelay.LogInfo = zaplog.Infof
	pushRtmprelay.ErrorInfo = zaplog.Errorf

	zaplog.Info("rtmprelay start push %s from %s", remoteurl, localurl)
	err := pushRtmprelay.StartWait(ctx)

	if err != nil {
		zaplog.Errorf("push error=%v", err)
		logger.Sync()
		os.Exit(-1)
	}
}
