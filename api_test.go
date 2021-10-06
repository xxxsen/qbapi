package qbapi

import (
	"context"
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
	if err := api.Login(context.Background()); err != nil {
		panic(err)
	}
	return api
}

func TestGetTorrentList(t *testing.T) {
	rsp, err := testApi.GetTorrentList(context.Background(), &GetTorrentListReq{})
	if err != nil {
		t.Fatal(err)
	}
	for _, item := range rsp.Items {
		t.Logf("data:%+v", *item)
	}
}

func TestGetApplicationVersion(t *testing.T) {
	rsp, err := testApi.GetApplicationVersion(context.Background(), &GetApplicationVersionReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("version:%+v", rsp)
}

func TestShutdownApplication(t *testing.T) {
	_, err := testApi.ShutDownAPPlication(context.Background(), &ShutdownApplicationReq{})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetApplicationPref(t *testing.T) {
	rsp, err := testApi.GetApplicationPreferences(context.Background(), &GetApplicationPreferencesReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetDefaultSavePath(t *testing.T) {
	rsp, err := testApi.GetDefaultSavePath(context.Background(), &GetDefaultSavePathReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetLog(t *testing.T) {
	rsp, err := testApi.GetLog(context.Background(), &GetLogReq{
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
	rsp, err := testApi.GetMainData(context.Background(), &GetMainDataReq{
		Rid: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetPeerHash(t *testing.T) {
	rsp, err := testApi.GetTorrentPeerData(context.Background(), &GetTorrentPeerDataReq{
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

func TestGetAlternativeSpeedLimitsState(t *testing.T) {
	rsp, err := testApi.GetAlternativeSpeedLimitsState(context.Background(), &GetAlternativeSpeedLimitsStateReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestToggleAlternativeSpeedLimits(t *testing.T) {
	rsp, err := testApi.ToggleAlternativeSpeedLimits(context.Background(), &ToggleAlternativeSpeedLimitsReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestGetGlobalDownloadLimit(t *testing.T) {
	rsp, err := testApi.GetGlobalDownloadLimit(context.Background(), &GetGlobalDownloadLimitReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestSetGlobalDownloadLimit(t *testing.T) {
	_, err := testApi.SetGlobalDownloadLimit(context.Background(), &SetGlobalDownloadLimitReq{
		Speed: 50 * 1024 * 1024, //50Mb/s
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetGlobalUploadLimit(t *testing.T) {
	rsp, err := testApi.GetGlobalUploadLimit(context.Background(), &GetGlobalUploadLimitReq{})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("data:%+v", rsp)
}

func TestSetGlobalUploadLimit(t *testing.T) {
	sp := 1.2 * 1024 * 1024
	_, err := testApi.SetGlobalUploadLimit(context.Background(), &SetGlobalUploadLimitReq{
		Speed: int(sp),
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBanPeers(t *testing.T) {
	_, err := testApi.BanPeers(context.Background(), &BanPeersReq{
		[]string{
			"54.111.178.247:18635",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}
