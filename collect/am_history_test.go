package collect

import "testing"

func TestQueryCrawlerRecord(t *testing.T) {

	reply,count,err:=QueryCrawlerRecord("Firefox for Android Beta",1,20)
	if err != nil {
		t.Errorf("TestQueryCrawlerRecord %v\n",err)
	}
	t.Log(reply,count)

}