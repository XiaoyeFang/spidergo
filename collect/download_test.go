package collect

import "testing"

func TestDownloadUpload(t *testing.T) {
	fid,key, err := DownloadUpload(
		//"com.chaozhuo.gameassistant",
		//"https://download.apkpure.com/b/apk/Y29tLmNoYW96aHVvLmdhbWVhc3Npc3RhbnRfMjEwXzIwZmQwYWI5?_fn=T2N0b3B1c192Mi4xLjBfYXBrcHVyZS5jb20uYXBr&k=8c16f791531b410206afef500ee43c4b5b35e807&as=6b684a2c5b2e18b08d4b2e78e3a9811d5b33457f&_p=Y29tLmNoYW96aHVvLmdhbWVhc3Npc3RhbnQ&c=1%7CTOOLS%7CZGV2PSVFNSU4QyU5NyVFNCVCQSVBQyVFOCVCNiU4NSVFNSU4RCU5MyVFNyVBNyU5MSVFNiU4QSU4MCVFNiU5QyU4OSVFOSU5OSU5MCVFNSU4NSVBQyVFNSU4RiVCOCVFRiVCQyU4OEJlaWppbmclMjBDaGFvemh1byUyMFRlY2hub2xvZ3klMjBDby4lMkMlMjBMdGQuJUVGJUJDJTg5JnZuPTIuMS4wJnZjPTIxMA"
		"Google App 8.10.17.21.arm64 beta (nodpi) (Android 5.0+)",
		"https://www.apkmirror.com/wp-content/themes/APKMirror/download.php?id=435802",

		)
	if err != nil {
		t.Error(err)
	}
	t.Logf("fid: %s,key: %v \n", fid,key)
}