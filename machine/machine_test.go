package machine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunServer_hello(t *testing.T) {
	server := NewServer()
	go server.Run(80)
	defer server.stop()
	resp, err := http.Get("http://localhost:80/")
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, "hello world", contents)
}
func TestRunServer_hello_8080_port(t *testing.T) {
	server := NewServer()
	go server.Run(8080)
	defer server.stop()
	resp, err := http.Get("http://localhost:8080/")
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, "hello world", contents)
}

// TODO : check the server log
func TestRunServer_hello_bad_port(t *testing.T) {
	server := NewServer()
	go server.Run(-1)
}

func TestAddHandlers(t *testing.T) {
	server := NewServer()
	server.AddHandlers([]Route{
		Route{"hello", "GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello world")
		}},
	})
	go server.Run(80)
	defer server.stop()
	resp, err := http.Get("http://localhost:80/hello")
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	contents, err := ioutil.ReadAll(resp.Body)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, "hello world", contents)
}
