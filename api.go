package qbapi

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/publicsuffix"
)

type QBAPI struct {
	c      *Config
	client *http.Client
}

func NewAPI(opts ...Option) (*QBAPI, error) {
	c := &Config{
		Timeout: 5 * time.Second,
	}
	for _, opt := range opts {
		opt(c)
	}
	if strings.HasSuffix(c.Host, "/") {
		c.Host = strings.TrimRight(c.Host, "/")
	}
	if len(c.Host) == 0 || len(c.Username) == 0 || len(c.Password) == 0 {
		return nil, NewMsgError(ErrParams, "params err")
	}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: c.Timeout,
		Jar:     jar,
		Transport: &http.Transport{
			DisableCompression: true,
		},
	}
	return &QBAPI{c: c, client: client}, nil
}

func (q *QBAPI) get(ctx context.Context, path string, req map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	if len(req) != 0 {
		query := httpReq.URL.Query()
		for k, v := range req {
			query.Add(k, v)
		}
		httpReq.URL.RawQuery = query.Encode()
	}

	resp, err := q.client.Do(httpReq)
	if err != nil {
		return nil, NewError(ErrNetwork, err)
	}
	return resp, nil
}

func (q *QBAPI) post(ctx context.Context, path string, values map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)

	var reader io.Reader
	if len(values) != 0 {
		form := url.Values{}
		for k, v := range values {
			form.Add(k, v)
		}
		reader = bytes.NewReader([]byte(form.Encode()))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, reader)
	if err != nil {
		return nil, NewError(ErrInternal, err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (q *QBAPI) struct2map(req interface{}) (map[string]string, error) {
	return ToMap(req, "json")
}

func (q *QBAPI) getWithDecoder(ctx context.Context, path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp, err := q.struct2map(req)
	if err != nil {
		return NewError(ErrMarsal, err)
	}
	httpRsp, err := q.get(ctx, path, mp)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return NewError(ErrStatusCode, fmt.Errorf("code:%d", httpRsp.StatusCode))
	}
	if rsp == nil {
		return nil
	}
	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return NewError(ErrNetwork, err)
	}
	if err := decoder(data, rsp); err != nil {
		return NewError(ErrUnmarsal, err)
	}
	return nil
}

func (q *QBAPI) buildURI(path string) string {
	return q.c.Host + path
}

func (q *QBAPI) postWithDecoder(ctx context.Context, path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp, err := q.struct2map(req)
	if err != nil {
		return NewError(ErrMarsal, err)
	}
	httpRsp, err := q.post(ctx, path, mp)
	if err != nil {
		return err
	}
	defer httpRsp.Body.Close()
	if httpRsp.StatusCode != http.StatusOK {
		return NewError(ErrStatusCode, fmt.Errorf("code:%d", httpRsp.StatusCode))
	}
	if rsp == nil {
		return nil
	}

	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return NewError(ErrNetwork, err)
	}
	if err := decoder(data, rsp); err != nil {
		return NewError(ErrUnmarsal, err)
	}
	return nil
}

//Login /api/v2/auth/login
func (q *QBAPI) Login(ctx context.Context) error {
	req := &LoginReq{Username: q.c.Username, Password: q.c.Password}
	var rsp string
	err := q.postWithDecoder(ctx, apiLogin, req, &rsp, StrDec)
	if err != nil {
		return err
	}
	if !strings.Contains(strings.ToLower(rsp), "ok") {
		return NewError(ErrLogin, fmt.Errorf("login fail:%s", rsp))
	}
	return nil
}

//GetApplicationVersion /api/v2/app/version
func (q *QBAPI) GetApplicationVersion(ctx context.Context, req *GetApplicationVersionReq) (*GetApplicationVersionRsp, error) {
	var version string
	err := q.getWithDecoder(ctx, apiGetAPPVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetApplicationVersionRsp{version}, nil
}

//GetAPIVersion /api/v2/app/webapiVersion
func (q *QBAPI) GetAPIVersion(ctx context.Context, req *GetAPIVersionReq) (*GetAPIVersionRsp, error) {
	var version string
	err := q.getWithDecoder(ctx, apiGetAPIVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetAPIVersionRsp{Version: version}, nil
}

//GetBuildInfo /api/v2/app/buildInfo
func (q *QBAPI) GetBuildInfo(ctx context.Context, req *GetBuildInfoReq) (*GetBuildInfoRsp, error) {
	rsp := &GetBuildInfoRsp{Info: &BuildInfo{}}
	if err := q.getWithDecoder(ctx, apiGetBuildInfo, req, rsp.Info, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//ShutDownAPPlication /api/v2/app/shutdown
func (q *QBAPI) ShutDownAPPlication(ctx context.Context, req *ShutdownApplicationReq) (*ShutdownApplicationRsp, error) {
	err := q.postWithDecoder(ctx, apiShutdownAPP, nil, nil, JsonDec)
	if err != nil {
		return nil, err
	}
	return &ShutdownApplicationRsp{}, nil
}

//GetApplicationPreferences /api/v2/app/preferences
func (q *QBAPI) GetApplicationPreferences(ctx context.Context, req *GetApplicationPreferencesReq) (*GetApplicationPreferencesRsp, error) {
	rsp := &GetApplicationPreferencesRsp{}
	err := q.getWithDecoder(ctx, apiGetAPPPerf, req, rsp, JsonDec)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetApplicationPreferences /api/v2/app/setPreferences
func (q *QBAPI) SetApplicationPreferences(ctx context.Context, req *SetApplicationPreferencesReq) (*SetApplicationPreferencesRsp, error) {
	//TODO:
	return nil, fmt.Errorf("not impl")
	// rsp := &SetApplicationPreferencesRsp{}
	// err := q.postWithDecoder(apiSetAPPPref, req, rsp, JsonDec)
	// if err != nil {
	// 	return nil, err
	// }
	// return rsp, nil
}

//GetDefaultSavePath /api/v2/app/defaultSavePath
func (q *QBAPI) GetDefaultSavePath(ctx context.Context, req *GetDefaultSavePathReq) (*GetDefaultSavePathRsp, error) {
	var path string
	if err := q.getWithDecoder(ctx, apiGetDefaultSavePath, nil, &path, StrDec); err != nil {
		return nil, err
	}
	return &GetDefaultSavePathRsp{Path: path}, nil
}

//GetLog /api/v2/log/main
func (q *QBAPI) GetLog(ctx context.Context, req *GetLogReq) (*GetLogRsp, error) {
	rsp := &GetLogRsp{Items: make([]*LogItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetPeerLog /api/v2/log/peers
func (q *QBAPI) GetPeerLog(ctx context.Context, req *GetPeerLogReq) (*GetPeerLogRsp, error) {
	rsp := &GetPeerLogRsp{Items: make([]*PeerLogItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetPeerLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetMainData /api/v2/sync/maindata
func (q *QBAPI) GetMainData(ctx context.Context, req *GetMainDataReq) (*GetMainDataRsp, error) {
	rsp := &GetMainDataRsp{}
	if err := q.getWithDecoder(ctx, apiGetMainData, req, &rsp, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentPeerData /api/v2/sync/torrentPeers
func (q *QBAPI) GetTorrentPeerData(ctx context.Context, req *GetTorrentPeerDataReq) (*GetTorrentPeerDataRsp, error) {
	rsp := &GetTorrentPeerDataRsp{Data: &TorrentPeerData{}, Exist: false}
	err := q.getWithDecoder(ctx, apiGetTorrentPeerData, req, rsp.Data, JsonDec)
	if err == nil {
		rsp.Exist = true
		return rsp, nil
	}
	code, err := RootCause(err)
	if code == ErrStatusCode && strings.Contains(err.Error(), "404") {
		return rsp, nil
	}
	return nil, err
}

//GetGlobalTransferInfo /api/v2/transfer/info
func (q *QBAPI) GetGlobalTransferInfo(ctx context.Context, req *GetGlobalTransferInfoReq) (*GetGlobalTransferInfoRsp, error) {
	rsp := &GetGlobalTransferInfoRsp{Info: &GlobalTransferInfo{}}
	if err := q.getWithDecoder(ctx, apiGetGlobalTransferInfo, req, rsp.Info, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetAlternativeSpeedLimitsState /api/v2/transfer/speedLimitsMode
func (q *QBAPI) GetAlternativeSpeedLimitsState(ctx context.Context, req *GetAlternativeSpeedLimitsStateReq) (*GetAlternativeSpeedLimitsStateRsp, error) {
	var intEnabled int
	rsp := &GetAlternativeSpeedLimitsStateRsp{Enabled: true}
	if err := q.getWithDecoder(ctx, apiGetAltSpeedLimitState, req, &intEnabled, IntDec); err != nil {
		return nil, err
	}
	if intEnabled == 0 {
		rsp.Enabled = false
	}
	return rsp, nil
}

//ToggleAlternativeSpeedLimits /api/v2/transfer/toggleSpeedLimitsMode
func (q *QBAPI) ToggleAlternativeSpeedLimits(ctx context.Context, req *ToggleAlternativeSpeedLimitsReq) (*ToggleAlternativeSpeedLimitsRsp, error) {
	rsp := &ToggleAlternativeSpeedLimitsRsp{}
	if err := q.postWithDecoder(ctx, apiToggleAltSpeedLimits, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetGlobalDownloadLimit /api/v2/transfer/downloadLimit
func (q *QBAPI) GetGlobalDownloadLimit(ctx context.Context, req *GetGlobalDownloadLimitReq) (*GetGlobalDownloadLimitRsp, error) {
	rsp := &GetGlobalDownloadLimitRsp{}
	if err := q.getWithDecoder(ctx, apiGetGlobalDownloadLimit, req, &rsp.Speed, IntDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetGlobalDownloadLimit /api/v2/transfer/setDownloadLimit
func (q *QBAPI) SetGlobalDownloadLimit(ctx context.Context, req *SetGlobalDownloadLimitReq) (*SetGlobalDownloadLimitRsp, error) {
	if err := q.postWithDecoder(ctx, apiSetGlobalDownloadLimit, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetGlobalDownloadLimitRsp{}, nil
}

//GetGlobalUploadLimit /api/v2/transfer/uploadLimit
func (q *QBAPI) GetGlobalUploadLimit(ctx context.Context, req *GetGlobalUploadLimitReq) (*GetGlobalUploadLimitRsp, error) {
	rsp := &GetGlobalUploadLimitRsp{}
	if err := q.getWithDecoder(ctx, apiGetGlobalUploadLimit, req, &rsp.Speed, IntDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetGlobalUploadLimit /api/v2/transfer/setUploadLimit
func (q *QBAPI) SetGlobalUploadLimit(ctx context.Context, req *SetGlobalUploadLimitReq) (*SetGlobalUploadLimitRsp, error) {
	if err := q.postWithDecoder(ctx, apiSetGlobalUploadLimit, req, nil, JsonDec); err != nil {
		return nil, err
	}
	return &SetGlobalUploadLimitRsp{}, nil
}

//BanPeers /api/v2/transfer/banPeers
func (q *QBAPI) BanPeers(ctx context.Context, req *BanPeersReq) (*BanPeersRsp, error) {
	for _, item := range req.Peers {
		if !strings.Contains(item, ":") {
			return nil, NewError(ErrParams, fmt.Errorf("invalid peer:%s", item))
		}
	}
	innerReq := &banPeersReqInner{Peers: strings.Join(req.Peers, "|")}
	if err := q.postWithDecoder(ctx, apiBanPeers, innerReq, nil, JsonDec); err != nil {
		return nil, err
	}
	return &BanPeersRsp{}, nil
}

//GetTorrentList /api/v2/torrents/info
func (q *QBAPI) GetTorrentList(ctx context.Context, req *GetTorrentListReq) (*GetTorrentListRsp, error) {
	rsp := &GetTorrentListRsp{Items: make([]*TorrentListItem, 0)}
	if err := q.getWithDecoder(ctx, apiGetTorrentList, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}
