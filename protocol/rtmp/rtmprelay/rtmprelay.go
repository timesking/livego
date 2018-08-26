package rtmprelay

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/timesking/livego/logs"
	"github.com/timesking/livego/protocol/amf"
	"github.com/timesking/livego/protocol/rtmp/core"
)

type RtmpRelay struct {
	PlayUrl              string
	PublishUrl           string
	cs_chan              chan core.ChunkStream
	stopChan             chan struct{}
	stopOnce             sync.Once
	connectPlayClient    *core.ConnClient
	connectPublishClient *core.ConnClient
	startflag            bool
	LastError            error
	DNSLookup            func(host string) (string, error)
}

func NewRtmpRelay(playurl *string, publishurl *string) *RtmpRelay {
	return &RtmpRelay{
		PlayUrl:              *playurl,
		PublishUrl:           *publishurl,
		cs_chan:              make(chan core.ChunkStream, 50),
		stopChan:             make(chan struct{}, 1),
		connectPlayClient:    nil,
		connectPublishClient: nil,
		startflag:            false,
	}
}

//TODO: ctrl-c not work when source connection is good
func (self *RtmpRelay) rcvPlayChunkStream(ctx context.Context) {
	defer self.Stop()
	logs.Debug("rcvPlayRtmpMediaPacket connectClient.Read START...")

loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-self.stopChan:
			break loop
		default:
		}
		var rc core.ChunkStream

		if self.startflag == false {
			self.connectPlayClient.Close(nil)
			logs.Info("rcvPlayChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
			break loop
		}
		self.connectPlayClient.SetReadDeadline(time.Now().Add(time.Second * 1))
		err := self.connectPlayClient.Read(&rc)

		// if err != nil && err == io.EOF {
		if err != nil {
			if err, ok := err.(net.Error); ok {
				if err.Timeout() {
					continue loop
				} else {
					self.LastError = err
					break loop
				}
			}
			if err == io.EOF {
				self.LastError = err
				break loop
			}
			logs.Error("rcvPlayChunkStream err: %s", err.Error())
		}
		//self.LogInfo("connectPlayClient.Read return rc.TypeID=%v length=%d, err=%v", rc.TypeID, len(rc.Data), err)
		switch rc.TypeID {
		case 20, 17:
			r := bytes.NewReader(rc.Data)
			vs, err := self.connectPlayClient.DecodeBatch(r, amf.AMF0)
			self.LastError = err
			logs.Error("rcvPlayRtmpMediaPacket: vs=%v, err=%v", vs, err)
			break
		case 18:
			logs.Debug("rcvPlayRtmpMediaPacket: metadata....")
		case 8, 9:
			self.cs_chan <- rc
		}
	}
	logs.Debug("rcvPlayRtmpMediaPacket connectClient.Read END...")
}

func (self *RtmpRelay) sendPublishChunkStream(ctx context.Context) {
	defer self.Stop()
	logs.Debug("sendPublishChunkStream START")
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case rc := <-self.cs_chan:
			//logs.Debug("sendPublishChunkStream: rc.TypeID=%v length=%d", rc.TypeID, len(rc.Data))
			self.connectPublishClient.SetReadDeadline(time.Now().Add(time.Second * 5))
			if err := self.connectPublishClient.Write(rc); err != nil {
				if err, ok := err.(net.Error); ok /*&& err.Timeout()*/ {
					self.LastError = err
					break loop
				}
				logs.Warn("sendPublishChunkStream err:", err.Error())
			}

		case <-self.stopChan:
			self.connectPublishClient.Close(nil)
			logs.Info("sendPublishChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
			break loop
		}
	}
	logs.Debug("sendPublishChunkStream END")
}

func (self *RtmpRelay) StartWait(ctx context.Context) error {
	if self.startflag {
		err := errors.New(fmt.Sprintf("The rtmprelay already started, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl))
		return err
	}

	self.connectPlayClient = core.NewConnClient()
	self.connectPublishClient = core.NewConnClient()

	logs.Debug("play server addr:%v starting....", self.PlayUrl)
	err := self.connectPlayClient.StartWithDNSLookup(self.PlayUrl, "play", nil)
	if err != nil {
		logs.Error("connectPlayClient.Start url=%v error", self.PlayUrl)
		return err
	}

	logs.Debug("publish server addr:%v starting....", self.PublishUrl)
	err = self.connectPublishClient.StartWithDNSLookup(self.PublishUrl, "publish", self.DNSLookup)
	if err != nil {
		logs.Error("connectPublishClient.Start url=%v error", self.PublishUrl)
		self.connectPlayClient.Close(nil)
		return err
	}

	self.startflag = true
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		self.rcvPlayChunkStream(ctx)
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		self.sendPublishChunkStream(ctx)
	}()

	wg.Wait()
	return self.LastError
}

func (self *RtmpRelay) Start() error {
	if self.startflag {
		err := errors.New(fmt.Sprintf("The rtmprelay already started, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl))
		return err
	}

	self.connectPlayClient = core.NewConnClient()
	self.connectPublishClient = core.NewConnClient()

	logs.Debug("play server addr:%v starting....", self.PlayUrl)
	err := self.connectPlayClient.Start(self.PlayUrl, "play")
	if err != nil {
		logs.Error("connectPlayClient.Start url=%v error", self.PlayUrl)
		return err
	}

	logs.Debug("publish server addr:%v starting....", self.PublishUrl)
	err = self.connectPublishClient.Start(self.PublishUrl, "publish")
	if err != nil {
		logs.Error("connectPublishClient.Start url=%v error", self.PublishUrl)
		self.connectPlayClient.Close(nil)
		return err
	}

	self.startflag = true
	ctx := context.TODO()
	go self.rcvPlayChunkStream(ctx)
	go self.sendPublishChunkStream(ctx)

	return nil
}

func (self *RtmpRelay) Stop() {
	self.stopOnce.Do(func() {
		if !self.startflag {
			logs.Debug("The rtmprelay already stoped, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
			return
		}
		self.startflag = false
		close(self.stopChan)
	})
}
