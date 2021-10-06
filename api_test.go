package qbapi

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

type cfg struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
}

var testApi = getAPI()

func getAPI() *QBAPI {
	data, err := ioutil.ReadFile(".vscode/cfg.json")
	if err != nil {
		panic(err)
	}
	cf := &cfg{}
	err = json.Unmarshal(data, cf)
	if err != nil {
		panic(err)
	}
	var opts []Option
	opts = append(opts, WithAuth(cf.Username, cf.Password))
	opts = append(opts, WithHost(cf.Host))
	api, err := NewAPI(opts...)
	if err != nil {
		panic(err)
	}
	if err := api.Login(); err != nil {
		panic(err)
	}
	return api
}

func TestGetTorrentList(t *testing.T) {
	rsp, err := testApi.GetTorrentList(&GetTorrentListRequest{})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetApplicationVersion(t *testing.T) {
	version, err := testApi.GetApplicationVersion()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("version:%s", version)
}
