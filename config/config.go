package config

import (
	"crawler-apkmirror/models"
	"encoding/json"
	"errors"
	iot "io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

const (
	MIRRPRURL        = "https://www.apkmirror.com"
	HTTP_USER_AGENT  = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.90 Safari/537.36"
	HTTP_ACCEPT      = "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8"
	URL_NOT_VALID    = "URL is not valid"
	EXISTS           = "Data already exists"
	MESSAGE          = "message"
	CONFIG           = "config"
	LASTUPDATE       = "Mirror the latest time to upload"
	NoUPDATEDVERSION = "No updated version"
	MAXBUCKETLIST    = 100
	ALLVERSIONS      = "All versions "
)

//AppConf global conf
var AppConf *models.AppConf

//TestData only debug
var TestData = []byte(`
name: crawler-am
http_port: 4003
grpc_listen: :5003
log_level: 10
proxy_http:
- http://172.16.8.8:31081
proxy_http_prefix: ''
check_update_interval: "@every 1h"
mongodb_url: mongodb://192.168.9.111:27017/apkpure_apkmirror
mongodb_name: "apkpure_apkmirror"
gp_detail_url: https://www.apkmirror.com
ip_test_url: http://v6.ip.zxinc.org/getip/
default_hl: en-US
s3_conf:
  endpoint: "http://ceph-rgw.staging.apkpure.com:80"
  access_key: "30P2Q3CJQAPASY9GBAON"
  secret_key: "qinXyPpDxgtGMlJQLoXgFajWu4PeCD1jlm7JXPVa"
  bucket: "apkmirror"
  prefix: "mirr"
  s3acl: "authenticated-read"
  chunk_size: 100MB
`)

func init() {
	//todo del. 目前方便测试
	var err error

	if AppConf == nil {
		AppConf, err = LoadConf("/conf/app.yml")
		if err != nil {
			panic(err)
		}
	}

}

func CreatDatabase() (db *mgo.Database, err error) {

	session, err := mgo.Dial(AppConf.MongodbUrl)
	if err != nil {

		panic(err)
	}
	db = session.DB(AppConf.MongodbName)

	return db, err
}

// LoadConf defines  how to load config from conf/app.yml or conf/app.json to AppConf.
func LoadConf(filepath string) (item *models.AppConf, err error) {
	if filepath == "" {
		return nil, errors.New("filepath is empty, must use --config xxx.yml/json")
	}

	data, err := iot.ReadFile(filepath)
	if err != nil {
		data = TestData
		glog.Infoln("debug mode,yaml is not use")
	}

	item = &models.AppConf{}
	if strings.HasSuffix(filepath, ".json") {
		err = json.Unmarshal(data, item)
	} else if strings.HasSuffix(filepath, ".yml") || strings.HasSuffix(filepath, ".yaml") {
		err = yaml.Unmarshal(data, item)
	} else {
		return nil, errors.New("you config file must be json/yml")
	}

	if err != nil {
		return nil, err
	}

	return item, nil
}

