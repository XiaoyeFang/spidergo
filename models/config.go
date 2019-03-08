package models

import "gopkg.in/mgo.v2/bson"

// AppConf defines  the value from configuration raw data.
type AppConf struct {
	Name                string   `yaml:"name" json:"name"`
	HTTPPort            int      `yaml:"http_port" json:"http_port"`
	GrpcListen          string   `yaml:"grpc_listen" json:"grpc_listen"`
	LogLevel            int      `yaml:"log_level" json:"log_level"`
	ProxyHttp           []string `yaml:"proxy_http" json:"proxy_http"`
	ProxyHttpPrefix     string   `yaml:"proxy_http_prefix" json:"proxy_http_prefix"`
	CheckUpdateInterval string   `yaml:"check_update_interval" json:"check_update_interval"`
	MongodbUrl          string   `yaml:"mongodb_url" json:"mongodb_url"`
	MongodbName         string   `yaml:"mongodb_name" json:"mongodb_name"`
	GpDetailUrl         string   `yaml:"gp_detail_url" json:"gp_detail_url"`
	GpSimilarUrl        string   `yaml:"gp_similar_url" json:"gp_similar_url"`
	GpSearchUrl         string   `yaml:"gp_search_url" json:"gp_search_url"`
	GpNewAppUrl         string   `yaml:"gp_new_app_url" json:"gp_new_app_url"`
	IpTestUrl           string   `yaml:"ip_test_url" json:"ip_test_url"`
	DefaultHl           string   `yaml:"default_hl" json:"default_hl"`
	DefaultGl           string   `yaml:"default_gl" json:"default_gl"`
	BurstLimit          int      `yaml:"burst_limit" json:"burst_limit"`
	S3Conf              S3Config `yaml:"s3_conf" json:"s3_conf"`
}

type S3Config struct {
	Endpoint  string `yaml:"endpoint" json:"endpoint"`
	AccessKey string `yaml:"access_key" json:"access_key"`
	SecretKey string `yaml:"secret_key" json:"secret_key"`
	Bucket    string `yaml:"bucket" json:"bucket"`
	Prefix    string `yaml:"prefix" json:"prefix"`
	S3Acl     string `yaml:"s3acl" json:"s3acl"`
	ChunkSize string `yaml:"chunk_size" json:"chunk_size"`
}

type HistoryList struct {
	// apkName: 查找的app name
	Id bson.ObjectId `json:"_id" bson:"_id"`
	// 要查询的apk名称，即用户传入的包名
	QueryName string `json:"query_name" bson:"query_name"`
	// apk 最新版本名称
	ApkName string `json:"apk_name" bson:"apk_name"`
	// apk 图标
	ApkIcon string `json:"apk_icon" bson:"apk_icon"`
	// apk 用户传入上传时间
	LastTime string `json:"last_time" bson:"last_time"`
	// apk 最新版本
	ApkVersion string `json:"apk_version" bson:"apk_version"`
	// apk 上传时间
	ApkUpload string `json:"apk_upload" bson:"apk_upload"`
	// apk 安装包大小
	ApkSize string `json:"apk_size" bson:"apk_size"`
	// apk 下载次数
	Downloads string `json:"downloads" bson:"downloads"`
	// apk 下载链接 同一个apk有不同的变种，会有多个url
	ApkLink []*DownloadLink `json:"apk_link" bson:"apk_link"`
	// apk 查询时间
	UpdateTime string `json:"update_time" bson:"update_time"`
}

type DownloadLink struct {
	// VariantName:包名
	VariantName string `json:"variant_name" bson:"variant_name"`
	// DownloadUrl: s3url
	DownloadUrl string `json:"download_url" bson:"download_url"`
	//MirrorUrl: apkmirror url
	MirrorUrl string `json:"mirror_url" bson:"mirror_url"`
	// code: 版本号
	VersionCode string `json:"version_code" bson:"version_code"`
}

type CollectDetail struct {
	// apkName: 查找的app name
	PackageName string `json:"pkg_name" bson:"pkg_name"`
	// url: 查找的apk
	Url string `json:"url" bson:"url"`
	// version: 版本
	LastUploadTime string `json:"last_upload_time" bson:"last_upload_time"`
	// time: 采集时间
	CrawlerTime string `json:"crawler_time" bson:"crawler_time"`
}
