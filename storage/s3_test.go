package storage

import (
	"strings"
	"testing"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"crawler-apkmirror/config"
	"crawler-apkmirror/models"
)

var Conf = config.AppConf.S3Conf

func init() {
	Conf = config.AppConf.S3Conf
}

//go test -v smart-backup/storage -run ^TestNewS3Storage$
func TestNewS3Storage(t *testing.T) {
	fs := NewS3Storage(&Conf)
	if fs == nil {
		t.Failed()
		return
	}

	fb := fs.Svc.Bucket(fs.Conf.Bucket)
	if fb == nil {
		t.Failed()
		return
	}

	//LsDir webp
	testDir := "test"
	keys, err := fs.LsDir(testDir)
	if err != nil {
		t.Error(err)
		return
	}
	for _, key := range keys {
		ret := strings.Split(key, "/")
		if len(ret) >= 2 && ret[1] != "" {
			t.Logf("key:%v", key)
		}
	}

	last := len(keys) - 1
	if last >= 1 {
		//rm
		t.Logf("will rm:%v", keys[last])
		exist, err := fb.Exists(keys[last])
		t.Logf("exist:%v err:%v", exist, err)
		if err != nil {
			t.Error(err)
			return
		}
		if exist {
			err = fs.RemoveDir(keys[1])
			if err != nil {
				t.Error(err)
				return
			}
		}
	}
	//upload
	err = fs.CopyFile("s3.go", "test/rm.go")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestS3Storage_CopyFile(t *testing.T) {
	t.Logf("Conf:%v", Conf.Endpoint)
	conf := models.S3Config{
		AccessKey: Conf.AccessKey,
		SecretKey: Conf.SecretKey,
		Endpoint:  Conf.Endpoint,
		Bucket:    Conf.Bucket,
		Prefix:    "prefix_youtube",
		ChunkSize: "100MB",
	}
	t.Logf("copyfile: conf:%v", conf)
	fs := NewS3Storage(&conf)
	if fs == nil {
		t.Failed()
		return
	}

	key := fmt.Sprintf("cike/2017-02-%d/s3.go", 10)
	fileList, err := fs.LsDir("cike")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("filelist:%v", fileList)
	err = fs.CopyFile("s3.go", key)
	if err != nil {
		t.Error(err)
		return
	}

}

//目前只支持1000个
func TestS3Storage_LsDir(t *testing.T) {
	conf := models.S3Config{
		AccessKey: Conf.AccessKey,
		SecretKey: Conf.SecretKey,
		Endpoint:  Conf.Endpoint,
		Bucket:    Conf.Bucket,
		Prefix:    "prefix_192_168_0_18",
		ChunkSize: "100MB",
	}
	fs := NewS3Storage(&conf)
	if fs == nil {
		t.Failed()
		return
	}

	testDir := "bigfile"
	keys, err := fs.LsDir(testDir)
	t.Logf("keys:%v", keys)
	if err != nil {
		t.Error(err)
		return
	}
	for _, key := range keys {
		if IsDir(key) {
			t.Logf("dir:%v", key)
		} else {
			t.Logf("file:%v", key)
		}
	}
}

func TestS3UploadDir(t *testing.T) {
	conf := models.S3Config{
		AccessKey: Conf.AccessKey,
		SecretKey: Conf.SecretKey,
		Endpoint:  Conf.Endpoint,
		Bucket:    Conf.Bucket,
		Prefix:    "imgcache",
		ChunkSize: "100MB",
	}
	fs := NewS3Storage(&conf)
	if fs == nil {
		t.Failed()
		return
	}

	backupFilePath := "/Users/apple/github/python/spider/data/imgcachedata"
	outPutPath := "/Users/apple/github/python/spider/data/imgcachedata"

	err := filepath.Walk(backupFilePath, func(localPath string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}

		if f.IsDir() {
			return nil
		}
		//path
		tmpKey := ""
		if strings.HasPrefix(localPath, outPutPath) {
			tmpKey = localPath[len(outPutPath):]
		}
		t := time.Now()
		y, m, d := t.Date()
		//key:pg/2017-02-02/a.txt
		key := fmt.Sprintf("/%4d-%02d-%02d%s", y, m, d, tmpKey)
		print("key:", key)
		//1.upload
		err = fs.CopyFile(localPath, key)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestGetFileType(t *testing.T) {
	tds := []struct {
		srcUrl   string
		fileType string
	}{
		{
			srcUrl:   "https://i.ytimg.com/vi/2EuTs1Yo-Bo/default.jpg",
			fileType: ".jpg",
		},
		{
			srcUrl:   "https://i.ytimg.com/vi/2EuTs1Yo-Bo/default.webp",
			fileType: ".webp",
		},
		{
			srcUrl:   "https://i.ytimg.com/vi/2EuTs1Yo-Bo/default.png",
			fileType: ".png",
		},
		{
			srcUrl:   "https://i.ytimg.com/vi/2EuTs1Yo-Bo/defaultpng",
			fileType: "",
		},
	}

	for _, td := range tds {
		fileType := GetFileType(td.srcUrl)
		if fileType != td.fileType {
			err := fmt.Sprintf("getFileType(%v)--->ret:%v want:%v", td.srcUrl, fileType, td.fileType)
			t.Fatal(err)
		}
		t.Logf("fileType:%v", fileType)
	}
}

func TestS3Storage_CheckExist(t *testing.T) {
	conf := models.S3Config{
		AccessKey: Conf.AccessKey,
		SecretKey: Conf.SecretKey,
		Endpoint:  Conf.Endpoint,
		Bucket:    Conf.Bucket,
		Prefix:    Conf.Prefix,
		ChunkSize: "100MB",
	}
	fs := NewS3Storage(&conf)
	if fs == nil {
		t.Failed()
		return
	}
	tds := []struct {
		key     string
		fid     string
		isExist bool
	}{
		{
			key:     "/F--SV-2A1x4",
			isExist: true,
		},
		{
			key:     "/F--SV-2A1x3",
			isExist: false,
		},
	}

	for index, td := range tds {
		fid, isExist := fs.CheckExist(td.key)
		if isExist != td.isExist {
			err := fmt.Errorf("index:%v fid:%v isExist:%v expect:%v key:%v", index, fid, isExist, td.isExist, td.key)
			t.Fatal(err)
		}
		t.Logf("fid:%v isExist:%v", fid, isExist)
	}
}
