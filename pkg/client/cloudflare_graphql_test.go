package client

import (
	"encoding/json"
	log "github.com/lastminutedotcom/heimdall/pkg/logging"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func Test_colocationDataCollection(t *testing.T) {
	file, _ := os.Open(filepath.Join("..", "..", "test", "cloudlfare_graphql_response.json"))
	defer file.Close()
	byteValue, _ := ioutil.ReadAll(file)
	response := model.Response{}
	err := json.Unmarshal(byteValue, &response)
	if err != nil {
		log.Info("An error occured: %v", err)
	}

	log.Info("response: %v", response)

}

func doHttpCall(request *http.Request) (*http.Response, error) {
	request = headers(request)
	return client.Do(request)
}

func headers(request *http.Request) *http.Request {
	for key, value := range fixHeaders() {
		request.Header.Set(key, value)
	}
	return request
}

func fixHeaders() map[string]string {
	return map[string]string{
		"X-Auth-Email": "api.sre@lastminute.com",
		"X-Auth-Key":   "3c77e8aaabe5377c8f5a037bbadcec230701b",
		"Content-Type": "application/json",
		"Accept":       "*/*",
	}
}
