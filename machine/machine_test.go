package machine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunServerHello(t *testing.T) {
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
func TestRunServerHello8080Port(t *testing.T) {
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
func TestRunServerHelloBadPort(t *testing.T) {
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
func TestAddHandlersCORS(t *testing.T) {
	server := NewServer()
	server.AddHandlers([]Route{
		Route{"hello", "GET", "/hello", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "hello world")
		}},
	})
	go server.Run(80)
	defer server.stop()
	req, err := http.NewRequest(http.MethodOptions, "http://localhost:80/hello", nil)
	req.Header.Add("Access-Control-Request-Method", "GET")
	req.Header.Add("Origin", "http://172.18.1.107:8081")
	req.Header.Add("Access-Control-Request-Headers", "X-Accesskeyid, X-Identityid, X-Secretkey, X-Sessiontoken")
	resp, err := http.DefaultClient.Do(req)
	assert.EqualValues(t, nil, err)
	assert.EqualValues(t, 200, resp.StatusCode)
	assert.EqualValues(t, "Origin", resp.Header.Get("Vary"))
	headers := resp.Header.Get("Access-Control-Allow-Headers")
	assert.EqualValues(t, "X-Accesskeyid, X-Identityid, X-Secretkey, X-Sessiontoken", headers, "Access-Control-Allow-Headers")
	method := resp.Header.Get("Access-Control-Allow-Methods")
	assert.EqualValues(t, "GET", method, "Access-Control-Allow-Methods")
	orrgin := resp.Header.Get("Access-Control-Allow-Origin")
	assert.EqualValues(t, "*", orrgin, "Access-Control-Allow-Origin")
}
