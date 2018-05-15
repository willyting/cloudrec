package machine

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestRunServer_hello(t *testing.T) {
	server := NewServer()
	go server.Run(80)
	defer server.stop()
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
func TestRunServer_hello_8080_port(t *testing.T) {
	server := NewServer()
	go server.Run(8080)
	defer server.stop()
	resp, err := http.Get("http://localhost:8080/")
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

// TODO : check the server log
func TestRunServer_hello_bad_port(t *testing.T) {
	server := NewServer()
	go server.Run(-1)
}

/*func TestAddHandlers(t *testing.T) {
	server := NewServer()
	server.AddHandlers([]Route{
		Route{"hello", "GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello world")
		}},
	})
	go server.Run(80)
	resp, err := http.Get("http://localhost:80/hello")
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
*/
