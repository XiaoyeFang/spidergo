package collect

import (
	"testing"
	"crawler-apkmirror/config"
	pb "crawler-apkmirror/protos"
)

func TestCheckUpdate(t *testing.T) {

	CheckUpdate("com.facebook.apk", "https://www.apkmirror.com/apk/facebook-2/facebook/", "179.0.0.36.82")


}
func TestGetAllVerLink(t *testing.T) {
	_, err := GetAllVerLink("/apk/mozilla/firefox-beta/firefox-beta-61-0-release/")
	if err != nil {
		t.Error(err)
	}

}

func TestGetDownloadLink(t *testing.T) {
	_, err := GetDownloadLink("https://www.apkmirror.com/apk/mozilla/firefox-beta/firefox-beta-61-0-release/firefox-android-beta-61-0-android-apk-download/")
	if err != nil {
		t.Error(err)
	}

}

func TestInsertApkRecord(t *testing.T) {
	db, err := config.CreatDatabase()
	if err != nil {
		t.Error(err)
	}
	reply := &pb.ApkDetailReply{}
	err = InsertApkRecord(db, reply, "Firefox for Android Beta", "60.0")

	if err != nil {
		t.Error(err)
	}

}

func TestUpdateFid(t *testing.T) {
	db, err := config.CreatDatabase()
	if err != nil {
		t.Error(err)
	}
	err = UpdateFid(db, "Google Sheets 1.18.252.06.33 (arm) (240dpi) (Android 5.0+)", "11111")

	if err != nil {
		t.Error(err)
	}
}

func TestInsertConfig(t *testing.T) {
	db, err := config.CreatDatabase()
	if err != nil {
		t.Error(err)
	}

	err =InsertConfig(db,"com.lawnchair.apk","https://www.apkmirror.com/apk/deletescape/lawnchair/","June 29, 2018 at 10:31PM GMT+0800")

	if err != nil {
		t.Error(err)
	}


}