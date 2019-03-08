package storage

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"crawler-apkmirror/config"
	"crawler-apkmirror/models"
	"crawler-apkmirror/ippoll"
	"github.com/golang/glog"

	"io"

	"github.com/dustin/go-humanize"
	"github.com/goamz/goamz/aws"
	"github.com/goamz/goamz/s3"
)

const FORMATS3PREFIX = "%s%s"
const FILECHUNKSIZE = 100 * 1024 * 1024 //bytes

type S3Storage struct {
	Svc  *s3.S3
	Conf models.S3Config
}

var extensionTypes = map[string]string{
	".apk":     "application/vnd.android.package-archive",
	".xapk":    "application/xapk-package-archive",
	".torrent": "application/x-bittorrent",
}
var s3AclValues = map[string]bool{
	"private":                   true,
	"public-read":               true,
	"public-read-write":         true,
	"authenticated-read":        true,
	"bucket-owner-full-control": true,
	"bucket-owner-read":         true,
}

func NewS3Storage(conf *models.S3Config) *S3Storage {
	if conf == nil {
		panic("NewS3Storage failed conf is nil")
	}
	fs := &S3Storage{}
	fs.Conf = *conf
	//check acl
	if _, ok := s3AclValues[conf.S3Acl]; !ok {
		fs.Conf.S3Acl = "private"
	}
	auth := aws.Auth{
		AccessKey: conf.AccessKey,
		SecretKey: conf.SecretKey,
	}

	fs.Svc = s3.New(auth, aws.Region{
		S3Endpoint:        conf.Endpoint,
		S3LowercaseBucket: true})

	if fs.Svc == nil {
		panic("NewS3Storage failed, fs.Svc is nil")
	}
	return fs
}

// TODO 实现 s3 storage
func (this *S3Storage) CopyFile(localPath, remotePath string) (err error) {
	//TODO remotePath: conf.prefix + remotepath
	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, _ := file.Stat()
	ext := strings.ToLower(path.Ext(localPath))
	var size int64 = fileInfo.Size()
	fileType := extensionTypes[ext]
	if fileType == "" {
		sniffLen := 512
		if size < int64(sniffLen) {
			sniffLen = int(size)
		}
		sniffBuffer := make([]byte, sniffLen)
		_, err = file.Read(sniffBuffer)
		if err != nil {
			return
		}
		_, err = file.Seek(0, 0)
		if err != nil {
			return
		}
		fileType = http.DetectContentType(sniffBuffer)
	}
	b := this.Svc.Bucket(this.Conf.Bucket)
	if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
		remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	}

	tmpFileChunkSize := uint64(0)
	tmpFileChunk, err := humanize.ParseBytes(this.Conf.ChunkSize)
	if err != nil {
		glog.Errorf("err:%v tmpFileChunk:%v", err, tmpFileChunk)
	} else {
		tmpFileChunkSize = uint64(tmpFileChunk)
	}

	if tmpFileChunkSize == 0 {
		tmpFileChunkSize = FILECHUNKSIZE
	}

	glog.Fatalf("remotepath:%s...size:%v filechunk:%v", remotePath, size, tmpFileChunkSize)
	if uint64(size) > tmpFileChunkSize {
		//大文件mul
		multi, err := b.InitMulti(remotePath, fileType, s3.ACL(this.Conf.S3Acl))
		if err != nil {
			err = fmt.Errorf("initMulti err:%v multi:%v b:%v remotepath:%v ", err, multi, b, remotePath)
			return err
		}
		parts, err := multi.PutAll(file, FILECHUNKSIZE)
		if err != nil {
			err = fmt.Errorf("initMulti err:%v multi:%v b:%v remotepath:%v ", err, multi, b, remotePath)
			return err
		}
		err = multi.Complete(parts)
	} else {
		//小文件，单个上传
		err = b.PutReader(remotePath, file, size, fileType, s3.ACL(this.Conf.S3Acl), s3.Options{})
	}
	return err
}

func (this *S3Storage) CopyFileByUrl(srcUrl, remotePath string) (fid string, err error) {
	if srcUrl == "" {
		return "", errors.New("S3Storage.CopyFileByUrl srcUrl is empty")
	}
	req, err := http.NewRequest(
		"GET",
		srcUrl,
		nil,
	)
	if err != nil {
		return "", err
	}
	glog.V(2).Infof("link:%v", srcUrl)
	req.Header.Set("Accept", config.HTTP_ACCEPT)
	req.Header.Set("User-Agent", config.HTTP_USER_AGENT)
	hc := ippoll.GetHc()
	resp, err := hc.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("getVideoInfo hc.Do url:%v resp.code:%v", srcUrl, resp.StatusCode)
		return "", err
	}
	defer resp.Body.Close()
	if resp.ContentLength <= 0 {
		err = fmt.Errorf("file size err:%v, src:%v", resp.ContentLength, srcUrl)
		return "", err
	}
	fileType := GetFileType(srcUrl)
	fid, err = this.CopyFileByReader(resp.Body, resp.ContentLength, fileType, remotePath)
	//b := this.Svc.Bucket(this.Conf.Bucket)
	//if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
	//	remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	//}
	////err = b.PutReader(remotePath, resp.Body, resp.ContentLength, resp.Header.Get("Content-Type"), s3.PublicReadWrite, s3.Options{})
	//err = b.PutReader(remotePath, resp.Body, resp.ContentLength, fileType, s3.ACL(this.Conf.S3Acl), s3.Options{})
	return fid, err
}

func (this *S3Storage) CopyFileByReader(r io.Reader, size int64, fileType, remotePath string) (fid string, err error) {
	b := this.Svc.Bucket(this.Conf.Bucket)
	if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
		remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	}
	err = b.PutReader(remotePath, r, size, fileType, s3.ACL(this.Conf.S3Acl), s3.Options{})
	//todo 是否增加bucket信息
	return remotePath, err
}

func (this *S3Storage) LsDir(remotePath string) ([]string, error) {
	remotePath = checkPath(remotePath)
	if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
		remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	}

	b := this.Svc.Bucket(this.Conf.Bucket)
	tmpFile, err := b.List(remotePath, "/", "", 1000)
	if err != nil {
		return nil, err
	}

	dstKeys := []string{}
	for _, key := range tmpFile.CommonPrefixes {
		if !strings.HasPrefix(key, remotePath) {
			continue
		}

		key = key[len(remotePath):]
		if key == "" {
			continue
		}
		dstKeys = append(dstKeys, key)
	}

	for _, key := range tmpFile.Contents {
		if !strings.HasPrefix(key.Key, remotePath) {
			continue
		}
		key.Key = key.Key[len(remotePath):]
		//正常s3不会返回空
		if key.Key == "" {
			continue
		}
		dstKeys = append(dstKeys, key.Key)
	}
	return dstKeys, nil
}

//返回keys
func (this *S3Storage) lsdir(remotePath, prefix, delim, marker string, max int) ([]string, error) {
	remotePath = checkPath(remotePath)
	if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
		remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	}

	b := this.Svc.Bucket(this.Conf.Bucket)
	tmpFile, err := b.List(remotePath, delim, marker, max)
	if err != nil {
		return nil, err
	}
	dstKeys := tmpFile.CommonPrefixes
	for _, key := range tmpFile.Contents {
		dstKeys = append(dstKeys, key.Key)
	}
	return dstKeys, nil
}
func (this *S3Storage) CheckExist(key string) (string, bool) {
	if !strings.HasPrefix(key, this.Conf.Prefix) {
		key = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, key)
	}

	b := this.Svc.Bucket(this.Conf.Bucket)
	resp, err := b.GetResponse(key)
	if err != nil {
		//todo 处理s3网络错误. 404才认为不存在.
		//s, ok := err.(*s3.Error)
		//if !ok{
		//}
		glog.Errorf("resp:%v err:%v key:%v", resp, err, key)
		return key, false
	}
	if resp.StatusCode == 200 {
		return key, true
	}
	return key, false
}

func (this *S3Storage) RemoveDir(remotePath string) (err error) {
	remotePath = checkPath(remotePath)

	if remotePath == "" {
		return errors.New("you will delete all file, not allow")
	}
	if !strings.HasPrefix(remotePath, this.Conf.Prefix) {
		remotePath = fmt.Sprintf(FORMATS3PREFIX, this.Conf.Prefix, remotePath)
	}
	//logger.Debugf("path:%v", remotePath)

	b := this.Svc.Bucket(this.Conf.Bucket)
	for {
		//对应path 前缀超过1000的key, 循环删除
		keys, err := this.lsdir(remotePath, "", "", "", 1000)
		if err != nil {
			return err
		}

		if len(keys) == 0 {
			break
		}

		for _, key := range keys {
			glog.Warningf("s3 del key :%v..b:%v", key, b)
			err = b.Del(key)
		}
	}

	return
}
func GetFileType(link string) string {
	if link == "" {
		return ""
	}

	ret := strings.Split(link, "/")
	if len(ret) < 2 {
		return ""
	}
	return path.Ext(ret[len(ret)-1])
}

func checkPath(remotePath string) string {
	if remotePath != "" && !strings.HasSuffix(remotePath, "/") {
		remotePath = remotePath + "/"
	}
	return remotePath
}

func IsDir(key string) bool {
	if key != "" && strings.HasSuffix(key, "/") {
		return true
	}
	return false
}
