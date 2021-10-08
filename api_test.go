package qbapi

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type cfg struct {
	Username    string   `json:"username"`
	Password    string   `json:"password"`
	Host        string   `json:"host"`
	TorrentFile []string `json:"torrent"`
	Link        []string `json:"link"`
	ValidHash   string   `json:"valid_hash"`
}

var testCfg = getCfg()
var testApi = getAPI()

func getCfg() *cfg {
	data, err := ioutil.ReadFile(".vscode/cfg.json")
	if err != nil {
		panic(err)
	}
	cf := &cfg{}
	err = json.Unmarshal(data, cf)
	if err != nil {
		panic(err)
	}
	return cf
}

func getAPI() *QBAPI {
	cf := testCfg
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
	limit := 10
	rsp, err := testApi.GetTorrentList(context.Background(), &GetTorrentListReq{
		Limit: &limit,
	})
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
	_, err := testApi.GetTorrentPeerData(context.Background(), &GetTorrentPeerDataReq{
		Hash: "1b175f0992fe932de8de33139698c1fd26988096",
		Rid:  0,
	})
	if err != nil {
		t.Fatal(err)
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

func TestAddTorrent(t *testing.T) {
	_, err := testApi.AddNewTorrent(context.Background(), &AddNewTorrentReq{
		File: testCfg.TorrentFile,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddLink(t *testing.T) {
	pause := true
	dlLimit := 512 * 1024
	_, err := testApi.AddNewLink(context.Background(), &AddNewLinkReq{
		Url: testCfg.Link,
		Meta: &AddTorrentMeta{
			Paused:  &pause,
			DlLimit: &dlLimit,
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBuildInfo(t *testing.T) {
	rsp, err := testApi.GetBuildInfo(context.Background(), &GetBuildInfoReq{})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp)
	assert.NotNil(t, rsp.Info)
	assert.NotEqual(t, "", rsp.Info.QT)
}

func TestGetTorrentGenericProperties(t *testing.T) {
	_, err := testApi.GetTorrentGenericProperties(context.Background(), &GetTorrentGenericPropertiesReq{Hash: "123"})
	assert.NoError(t, err)
	rsp, err := testApi.GetTorrentGenericProperties(context.Background(), &GetTorrentGenericPropertiesReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Property)
}

func TestGetTorrentTrackers(t *testing.T) {
	rsp, err := testApi.GetTorrentTrackers(context.Background(), &GetTorrentTrackersReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	assert.NotEqual(t, 0, len(rsp.Trackers))
	t.Logf("data:%+v", rsp.Trackers)
}

func TestGetTorrentWebSeeds(t *testing.T) {
	rsp, err := testApi.GetTorrentWebSeeds(context.Background(), &GetTorrentWebSeedsReq{
		Hash: testCfg.ValidHash,
	})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.WebSeeds)
}

func TestGetTorrentContents(t *testing.T) {
	rsp, err := testApi.GetTorrentContents(context.Background(), &GetTorrentContentsReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.Contents)
}

func TestGetTorrentPiecesStates(t *testing.T) {
	rsp, err := testApi.GetTorrentPiecesStates(context.Background(), &GetTorrentPiecesStatesReq{Hash: testCfg.ValidHash})
	assert.NoError(t, err)
	t.Logf("data:%+v", rsp.States)
}
