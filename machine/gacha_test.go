package machine

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRunServer_hello(t *testing.T) {
	go RunServer(80)
	resp, err := http.Get("http://localhost:80/")
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != 200 {
		t.Error("get a error status")
	}
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if len(contents) <= 0 {
		t.Error("no response")
	}
	if string(contents) != "hello world" {
		t.Error("response incorrect")
	}
}
