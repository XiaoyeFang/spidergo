package collect

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"crawler-apkmirror/config"
	"crawler-apkmirror/ippoll"
	"crawler-apkmirror/models"
	pb "crawler-apkmirror/protos"
	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
	"github.com/robfig/cron"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func CheckUpdate(name, url, lastTime string) {
	c := cron.New()
	db, err := config.CreatDatabase()
	if err != nil {

		glog.V(0).Infof("CheckUpdate CreatDatabase err %v\n", err)
		panic(err)

	}

	defer func() {

		if r := recover(); r != nil {

			fmt.Println("Recovered in f", r)
			c.AddFunc(config.AppConf.CheckUpdateInterval, func() {
				TimerTask(db, name, url, lastTime)
			})
			c.Start()
			c.Run()
		}
	}()

	TimerTask(db, name, url, lastTime)

	//定时任务
	c.AddFunc(config.AppConf.CheckUpdateInterval, func() {

		TimerTask(db, name, url, lastTime)

	})

	c.Start()
	c.Run()

}

func TimerTask(db *mgo.Database, name, url, lastTime string) {

	reply := &pb.ApkDetailReply{}
	req, err := MakeCheckUrl(url)
	if err != nil {
		glog.Errorf("MakeCheckUrl %s\n", err)
		return
	}

	hc := ippoll.GetHc()
	resp, err := hc.Do(req)
	if err != nil {
		glog.Errorf("hc.Do(req) %s\n", err)
		return
	}

	if resp.StatusCode != 200 {
		fmt.Printf("TimerTask hc.Do url:%v resp.code:%v ip:%v", req.URL.RequestURI(), resp.StatusCode, hc.ProxyAddr)
		fmt.Println(resp.Header.Get("Retry-After"))
		return
	}
	reply, err = ParseAllVerHtml(resp)
	if err != nil {
		glog.V(0).Infof("ParseAllVerHtml %v \n", err)
		return
	}
	qureyTime, _ := time.Parse("2006-01-02 15:04:05", lastTime)
	UpdateTime, _ := time.Parse("2006-01-02 15:04:05", reply.ApkUpload)

	if UpdateTime.Before(qureyTime) {
		//无新版本
		glog.V(0).Infoln(name, config.LASTUPDATE, reply.ApkUpload, config.NoUPDATEDVERSION)
		return
	}
	//有新版本，下载保存，保存记录到mongodb,downloadurl改为fid
	err = InsertApkRecord(db, reply, name, lastTime)

	if err != nil {
		glog.V(0).Infof("InsertApkRecord err %v\n", err)
	}

	if len(reply.ApkLink) != 0 {

		//测试不开
		for _, v := range reply.ApkLink {
			//fmt.Println("DownloadUrl", v.DownloadUrl)
			fid, key, err := DownloadUpload(v.VariantName, v.DownloadUrl)
			if err != nil {
				glog.V(0).Infof("DownloadUpload err %s\n", err)
			}
			if fid == "" {
				fid = config.AppConf.S3Conf.Bucket + "/" + config.AppConf.S3Conf.Prefix + key
			}
			err = UpdateFid(db, v.VariantName, fid)
			if err != nil {
				glog.V(0).Infof("UpdateFid %v \n", err)
			}

		}
	}

}

func GetAllVerLink(fileLink string) (reply []*pb.DownloadLink, err error) {
	//fmt.Println("GetAllVerLink====", fileLink)
	req, err := MakeCheckUrl(config.MIRRPRURL + fileLink)

	if err != nil {
		return reply, err
	}
	hc := ippoll.GetHc()
	resp, err := hc.Do(req)

	if err != nil {
		return reply, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("Search hc.Do url:%v resp.code:%v ip:%v \n", req.URL.RequestURI(), resp.StatusCode, hc.ProxyAddr)
		fmt.Printf("请在 %s s后尝试请求 \n", resp.Header.Get("Retry-After"))

		return nil, err
	}
	reply, err = ParseLatestVerHtml(resp)

	return reply, err
}

func GetLatestVerLink(fileLink string) (reply *pb.DownloadLink, err error) {
	//fmt.Println("GetLatestVerLink==", fileLink)
	req, err := MakeCheckUrl(fileLink)
	if err != nil {
		glog.Infoln(err)
		return reply, err
	}
	hc := ippoll.GetHc()
	resp, err := hc.Do(req)
	if err != nil {
		return reply, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetDownloadLink hc.Do url:%v resp.code:%v ip:%v", req.URL.RequestURI(), resp.StatusCode, hc.ProxyAddr)
		return nil, err
	}
	reply, err = ParseDownLoadHtml(resp)
	if err != nil {
		return nil, err
	}

	return reply, err
}

func GetDownloadLink(downloadUrl string) (url string, err error) {
	//fmt.Println("downloadUrl===",downloadUrl)
	req, err := MakeCheckUrl(downloadUrl)
	if err != nil {
		glog.Infoln(err)
		return url, err
	}
	hc := ippoll.GetHc()

	resp, err := hc.Do(req)
	if err != nil {
		return url, err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("GetDownloadLink hc.Do url:%v resp.code:%v ip:%v", req.URL.RequestURI(), resp.StatusCode, hc.ProxyAddr)
		return url, err
	}

	downloadDoc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		glog.V(0).Infoln("downloadDoc", err)
		return "", err
	}

	url, _ = downloadDoc.Find(".noPadding .notes a").Attr("href")
	//text := downloadDoc.Find(".appspec-value").Text()
	//fmt.Println("text==========", text)
	//
	//if text != "" {
	//	versionCode = text[strings.Index(text, "(")+1:strings.Index(text, ")")]
	//}

	return url, err
}

func InsertApkRecord(db *mgo.Database, reply *pb.ApkDetailReply, apkName, lastTime string) error {

	history := &models.HistoryList{}
	history.ApkLink = make([]*models.DownloadLink, 0, 0)
	msgCol := db.C(config.MESSAGE)
	count, err := msgCol.Find(bson.M{"query_name": apkName, "last_time": lastTime, "apk_upload": reply.ApkUpload}).Count()

	if err != nil {
		glog.V(0).Infof("msgCol.Find %s\n", err)

		return err
	}

	if count != 0 {
		err = fmt.Errorf("%s", config.EXISTS)

		return err

	}
	glog.V(0).Infof("msgCol.Find Count %d\n", count)
	history.Id = bson.NewObjectId()
	history.QueryName = apkName
	history.ApkName = reply.ApkName
	history.LastTime = lastTime
	history.ApkVersion = reply.ApkVersion
	history.ApkSize = reply.ApkSize
	history.ApkUpload = reply.ApkUpload
	history.ApkIcon = reply.ApkIcon
	history.Downloads = reply.Downloads
	history.UpdateTime = time.Now().String()

	if len(reply.ApkLink) != 0 {

		for _, v := range reply.ApkLink {
			apkLink := &models.DownloadLink{}
			apkLink.VariantName = v.VariantName
			//mirrurl
			apkLink.MirrorUrl = v.DownloadUrl
			DownloadUrl := fmt.Sprintf("%s03d", v.DownloadUrl)
			DownloadUrl = base64.RawURLEncoding.EncodeToString([]byte(DownloadUrl))
			//s3url
			apkLink.DownloadUrl = config.AppConf.S3Conf.Bucket + "/" + config.AppConf.S3Conf.Prefix + DownloadUrl
			apkLink.VersionCode = v.VersionCode
			history.ApkLink = append(history.ApkLink, apkLink)
		}
	}

	err = msgCol.Insert(history)

	if err != nil {
		glog.V(0).Infof("msgCol.Insert %s\n", err)
		return err
	}

	return err
}

func UpdateFid(db *mgo.Database, variantName, fid string) (err error) {

	dcol := db.C(config.MESSAGE)

	selector := bson.M{"apk_link.variant_name": variantName}

	data := bson.M{"$set": bson.M{"apk_link.$.download_url": fid}}

	err = dcol.Update(selector, data)

	if err != nil {

		return err
	}
	return err
}

func InsertConfig(db *mgo.Database, packageName, url, lastUploadTime string) (err error) {

	dcol := db.C(config.CONFIG)
	if !strings.HasPrefix(url, "http") {
		err = fmt.Errorf("url %s", config.URL_NOT_VALID)
		return err
	}

	count, err := dcol.Find(bson.M{"pkg_name": packageName, "url": url}).Count()
	if err != nil {
		glog.V(0).Infof("InsertConfig err%v \n", err)
	}
	//fmt.Println(count)
	if count != 0 {

		//已存在该包配置则更新lastTime以及url

		selector := bson.M{"pkg_name": packageName}

		data := bson.M{"$set": bson.M{"url": url, "last_upload_time": lastUploadTime, "crawler_time": time.Now().String()}}

		err = dcol.Update(selector, data)

		if err != nil {

			return err
		}
		return err
	}

	details := models.CollectDetail{
		packageName,
		url,
		lastUploadTime,
		time.Now().String(),}

	//fmt.Println(details)

	err = dcol.Insert(details)
	if err != nil {
		glog.Errorf("dcol Insert err %v\n", err)
		return err
	}

	return err
}
