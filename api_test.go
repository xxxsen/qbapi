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
	rsp, err := testApi.GetTorrentList(&GetTorrentListReq{})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetApplicationVersion(t *testing.T) {
	rsp, err := testApi.GetApplicationVersion(&GetApplicationVersionReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("version:%+v", rsp)
}

func TestShutdownApplication(t *testing.T) {
	_, err := testApi.ShutDownAPPlication(&ShutdownApplicationReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetApplicationPref(t *testing.T) {
	rsp, err := testApi.GetApplicationPreferences(&GetApplicationPreferencesReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetDefaultSavePath(t *testing.T) {
	rsp, err := testApi.GetDefaultSavePath(&GetDefaultSavePathReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetLog(t *testing.T) {
	rsp, err := testApi.GetLog(&GetLogReq{
		Normal:   true,
		Info:     true,
		Warning:  true,
		Critical: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetMainData(t *testing.T) {
	rsp, err := testApi.GetMainData(&GetMainDataReq{
		Rid: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetPeerHash(t *testing.T) {
	rsp, err := testApi.GetTorrentPeerData(&GetTorrentPeerDataReq{
		Hash: "1b175f0992fe932de8de33139698c1fd26988096",
		Rid:  0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if rsp.Exist {
		t.Logf("found:%+v", rsp.Data)
	} else {
		t.Log("not found")
	}
}
