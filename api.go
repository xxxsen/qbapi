package qbapi

import (
	"bytes"
	"encoding/json"
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

func (q *QBAPI) get(path string, req map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)
	httpReq, err := http.NewRequest(http.MethodGet, uri, nil)
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

func (q *QBAPI) post(path string, values map[string]string) (*http.Response, error) {
	uri := q.buildURI(path)

	var reader io.Reader
	if len(values) != 0 {
		form := url.Values{}
		for k, v := range values {
			form.Add(k, v)
		}
		reader = bytes.NewReader([]byte(form.Encode()))
	}

	req, err := http.NewRequest(http.MethodPost, uri, reader)
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

func (q *QBAPI) struct2map(req interface{}) map[string]string {
	if req == nil {
		return nil
	}
	mp := make(map[string]string)
	data, _ := json.Marshal(req)
	_ = json.Unmarshal(data, &mp)
	return mp
}

func (q *QBAPI) getWithDecoder(path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp := q.struct2map(req)
	httpRsp, err := q.get(path, mp)
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

func (q *QBAPI) postWithDecoder(path string, req interface{}, rsp interface{}, decoder Decoder) error {
	mp := q.struct2map(req)
	httpRsp, err := q.post(path, mp)
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
func (q *QBAPI) Login() error {
	req := &LoginReq{Username: q.c.Username, Password: q.c.Password}
	rsp, err := q.post(apiLogin, q.struct2map(req))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return NewError(ErrNetwork, err)
	}
	if !strings.Contains(strings.ToLower(string(data)), "ok") {
		return NewError(ErrLogin, fmt.Errorf("login fail:%s", string(data)))
	}
	return nil
}

//GetApplicationVersion /api/v2/app/version
func (q *QBAPI) GetApplicationVersion(req *GetApplicationVersionReq) (*GetApplicationVersionRsp, error) {
	var version string
	err := q.getWithDecoder(apiGetAPPVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetApplicationVersionRsp{version}, nil
}

//GetAPIVersion /api/v2/app/webapiVersion
func (q *QBAPI) GetAPIVersion(req *GetAPIVersionReq) (*GetAPIVersionRsp, error) {
	var version string
	err := q.getWithDecoder(apiGetAPIVersion, nil, &version, StrDec)
	if err != nil {
		return nil, err
	}
	return &GetAPIVersionRsp{Version: version}, nil
}

//GetBuildInfo /api/v2/app/buildInfo
func (q *QBAPI) GetBuildInfo(req *GetBuildInfoReq) (*GetBuildInfoRsp, error) {
	rsp := &GetBuildInfoRsp{Info: &BuildInfo{}}
	if err := q.getWithDecoder(apiGetBuildInfo, req, rsp.Info, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//ShutDownAPPlication /api/v2/app/shutdown
func (q *QBAPI) ShutDownAPPlication(req *ShutdownApplicationReq) (*ShutdownApplicationRsp, error) {
	err := q.postWithDecoder(apiShutdownAPP, nil, nil, JsonDec)
	if err != nil {
		return nil, err
	}
	return &ShutdownApplicationRsp{}, nil
}

//GetApplicationPreferences /api/v2/app/preferences
func (q *QBAPI) GetApplicationPreferences(req *GetApplicationPreferencesReq) (*GetApplicationPreferencesRsp, error) {
	rsp := &GetApplicationPreferencesRsp{}
	err := q.getWithDecoder(apiGetAPPPerf, req, rsp, JsonDec)
	if err != nil {
		return nil, err
	}
	return rsp, nil
}

//SetApplicationPreferences /api/v2/app/setPreferences
func (q *QBAPI) SetApplicationPreferences(req *SetApplicationPreferencesReq) (*SetApplicationPreferencesRsp, error) {
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
func (q *QBAPI) GetDefaultSavePath(req *GetDefaultSavePathReq) (*GetDefaultSavePathRsp, error) {
	var path string
	if err := q.getWithDecoder(apiGetDefaultSavePath, nil, &path, StrDec); err != nil {
		return nil, err
	}
	return &GetDefaultSavePathRsp{Path: path}, nil
}

//GetLog /api/v2/log/main
func (q *QBAPI) GetLog(req *GetLogReq) (*GetLogRsp, error) {
	rsp := &GetLogRsp{Items: make([]*LogItem, 0)}
	if err := q.getWithDecoder(apiGetLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetPeerLog /api/v2/log/peers
func (q *QBAPI) GetPeerLog(req *GetPeerLogReq) (*GetPeerLogRsp, error) {
	rsp := &GetPeerLogRsp{Items: make([]*PeerLogItem, 0)}
	if err := q.getWithDecoder(apiGetPeerLog, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetMainData /api/v2/sync/maindata
func (q *QBAPI) GetMainData(req *GetMainDataReq) (*GetMainDataRsp, error) {
	rsp := &GetMainDataRsp{}
	if err := q.getWithDecoder(apiGetMainData, req, &rsp, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}

//GetTorrentPeerData /api/v2/sync/torrentPeers
func (q *QBAPI) GetTorrentPeerData(req *GetTorrentPeerDataReq) (*GetTorrentPeerDataRsp, error) {
	rsp := &GetTorrentPeerDataRsp{Data: &TorrentPeerData{}, Exist: false}
	err := q.getWithDecoder(apiGetTorrentPeerData, req, rsp.Data, JsonDec)
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

//GetTorrentList /api/v2/torrents/info
func (q *QBAPI) GetTorrentList(req *GetTorrentListReq) (*GetTorrentListRsp, error) {
	rsp := &GetTorrentListRsp{Items: make([]*TorrentListItem, 0)}
	if err := q.getWithDecoder(apiGetTorrentList, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}
