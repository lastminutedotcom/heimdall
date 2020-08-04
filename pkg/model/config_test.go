package model_test

import (
	"bytes"
	"github.com/lastminutedotcom/heimdall/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var sampleConf = `
{
  "collect_every_minutes" : "5",
  "graphite_config": {
    "host": "graphite.company.com",
    "port": 2113
  },
  "kubernetes": {
    "management_port": "8888"
  }
}
`

func Test_ParseConfig(t *testing.T){
	cfg, err := model.ParseConfig(bytes.NewBufferString(sampleConf))
	if err!=nil {
		t.Fatalf("could not read config: %v", err)
	}
	assert.Equal(t, cfg.CollectEveryMinutes, "5")
	assert.Equal(t, cfg.KubeConfig.MgmtPort, "8888")
}
