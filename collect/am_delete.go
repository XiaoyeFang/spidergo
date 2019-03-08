package collect

import (
	"crawler-apkmirror/config"
	"crawler-apkmirror/storage"
	"fmt"
)


func DeleteBucket() (err error) {

	fs := storage.NewS3Storage(&config.AppConf.S3Conf)

	err = fs.RemoveDir(config.AppConf.S3Conf.Prefix)
	if err != nil {

		fmt.Printf("fs.RemoveDir err %v\n", err)
		return err
	}
	fmt.Println("Delete Success")

	//db, err := config.CreatDatabase()
	//
	//dcol := db.C(config.MESSAGE)
	//
	//err = dcol.Remove(bson.M{"query_name": "com.lawnchair.apk"})

	//err = dcol.Find(bson.M{}).All(&objects.Objects)
	//if err != nil {
	//	glog.V(0).Infof("dcol.Find %s\n", err)
	//
	//	return err
	//}

	return err
}
