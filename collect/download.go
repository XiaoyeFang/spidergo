package collect

import (
	"fmt"
	"encoding/base64"
	"crawler-apkmirror/storage"
	"crawler-apkmirror/config"
	"github.com/golang/glog"
)

//下载apk
func DownloadUpload(packageName, link string) (fid, key string, err error) {
	//fmt.Printf("\n %c[1;40;32m packageName  %s link  %s %c[0m\n\n", 0x1B, packageName, link, 0x1B)
	fs := storage.NewS3Storage(&config.AppConf.S3Conf)
	//key = fmt.Sprintf("%s_%d_%03d", packageName, time.Now().Unix(), rand.Intn(100))
	key = fmt.Sprintf("%s03d", link)
	//base64
	key = base64.RawURLEncoding.EncodeToString([]byte(key))
	fid, err = fs.CopyFileByUrl(link, key)
	if err != nil {
		fmt.Errorf("\n %c[1;40;32m %v %c[0m\n\n", 0x1B, err, 0x1B)
		return "", key, err
	}
	if fs.Conf.Bucket == "" {
		//never do
		glog.V(0).Infof("fs err:%v", fs)
	}
	fid = fmt.Sprintf("%s/%s", fs.Conf.Bucket, fid)
	fmt.Printf("\n %c[1;40;32m packageName  %s fid  %s%c[0m\n\n", 0x1B, packageName, fid, 0x1B)

	return fid, key, err
}
