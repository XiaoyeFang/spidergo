package collect

import (
	"crawler-apkmirror/config"
	"crawler-apkmirror/models"
	pb "crawler-apkmirror/protos"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2/bson"
)

func QueryCrawlerRecord(apkName string, page, pageSize int32) (reply []*pb.ApkDetailReply, count int, err error) {
	var queryReply = make([]models.HistoryList, 0)
	reply = make([]*pb.ApkDetailReply, 0)
	p := int(page)
	ps := int(pageSize)
	if ps == 0 {
		ps = 20
	}

	db, err := config.CreatDatabase()
	if err != nil {
		glog.V(0).Infof("config.CreatDatabase %s\n", err)
		return reply, count, err
	}
	msgCol := db.C(config.MESSAGE)

	if apkName == "" {
		err = msgCol.Find(bson.M{}).Skip((p - 1) * ps).Limit(ps).All(&queryReply)
		count, err = msgCol.Find(bson.M{}).Count()

	} else {
		err = msgCol.Find(bson.M{"queryname": apkName}).Skip((p - 1) * ps).Limit(ps).All(&queryReply)
		count, err = msgCol.Find(bson.M{"queryname": apkName}).Count()
	}

	if err != nil {
		glog.V(0).Infof("msgCol.Find err%v\n", err)
		return nil, count, err
	}
	if len(queryReply) != 0 {
		for _, v := range queryReply {
			detail := &pb.ApkDetailReply{}
			detail.Id = v.Id.Hex()
			detail.ApkName = v.ApkName
			detail.QueryName = v.QueryName
			detail.ApkIcon = v.ApkIcon
			detail.Downloads = v.Downloads
			detail.ApkVersion = v.ApkVersion
			detail.ApkUpload = v.ApkUpload
			detail.ApkSize = v.ApkSize
			detail.UpdateTime = v.UpdateTime

			if len(v.ApkLink) != 0 {
				for _, k := range v.ApkLink {
					apkLink := &pb.DownloadLink{}
					apkLink.VariantName = k.VariantName
					apkLink.DownloadUrl = k.DownloadUrl
					detail.ApkLink = append(detail.ApkLink, apkLink)
				}

			}

			reply = append(reply, detail)

		}
	}

	//glog.V(0).Infof("reply %v\n", reply)
	return reply, count, err
}
