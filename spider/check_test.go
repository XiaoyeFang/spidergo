package spider

import "testing"

func TestCheckAPKUpdate(t *testing.T) {
	CheckAPKUpdate("com.facebook.apk", "https://www.apkmirror.com/apk/google-inc/google-play-store/google-play-store-12-6-13-release/", "179.0.0.36.82")

}
