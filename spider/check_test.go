package spider

import "testing"

func TestCheckAPKUpdate(t *testing.T) {
	CheckAPKUpdate("com.facebook.apk", "https://www.apkmirror.com/apk/facebook-2/facebook/", "179.0.0.36.82")

}
