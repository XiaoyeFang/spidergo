package collect

import (
	"golang.org/x/net/context"
	pb "crawler-apkmirror/protos"
	"github.com/golang/glog"
	"crawler-apkmirror/config"
	"fmt"
)

type CollectServer struct {
}

func (*CollectServer) CheckUpdate(ctx context.Context, req *pb.CrawlerConfigRequest) (*pb.CrawlerConfigReply, error) {
	reply := &pb.CrawlerConfigReply{}
	//保存配置到mongodb
	db, err := config.CreatDatabase()
	if err != nil {
		panic(err)
	}
	if len(req.CheckList) != 0 {
		for _, config := range req.CheckList {

			err := InsertConfig(db, config.PackageName, config.Url, config.LastUploadTime)
			if err != nil {

			}
			go CheckUpdate(config.PackageName, config.Url, config.LastUploadTime)

		}
	}

	return reply, err

}

func (*CollectServer) CollectHistory(ctx context.Context, req *pb.CrawlerTaskRequest) (*pb.CrawlerTaskReply, error) {
	fmt.Println("====CollectHistory====",req)
	reply := &pb.CrawlerTaskReply{}

	appInfoList, count, err := QueryCrawlerRecord(req.PackageName, req.Page, req.PageSize)

	if err != nil {
		glog.V(0).Infof("CollectHistory err %v\n", err)

		return reply, err
	}
	reply.AppInfoList, reply.Count = appInfoList, int32(count)
	//fmt.Println(count)
	return reply, nil
}
