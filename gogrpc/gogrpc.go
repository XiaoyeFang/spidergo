package gogrpc

import (
	"github.com/XiaoyeFang/spidergo/protos"
	"golang.org/x/net/context"
)

type Spidergo struct {
}

func (this *Spidergo) CheckUpdate(ctx context.Context, request protos.CrawlerConfigRequest) (reply protos.CrawlerConfigReply, err error) {
	return protos.CrawlerConfigReply{}, nil
}

