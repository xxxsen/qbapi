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
		return nil, fmt.Errorf("params err")
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
		return nil, err
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
		return nil, err
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
		return nil, err
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
		return fmt.Errorf("code:%d", httpRsp.StatusCode)
	}
	if rsp == nil {
		return nil
	}
	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return err
	}
	if err := decoder(data, rsp); err != nil {
		return err
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
		return fmt.Errorf("code:%d", httpRsp.StatusCode)
	}
	if rsp == nil {
		return nil
	}

	data, err := ioutil.ReadAll(httpRsp.Body)
	if err != nil {
		return err
	}
	if err := decoder(data, rsp); err != nil {
		return err
	}
	return nil
}

func (q *QBAPI) Login() error {
	req := &LoginReq{Username: q.c.Username, Password: q.c.Password}
	rsp, err := q.post(apiLogin, q.struct2map(req))
	if err != nil {
		return err
	}
	defer rsp.Body.Close()
	data, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return err
	}
	if !strings.Contains(strings.ToLower(string(data)), "ok") {
		return fmt.Errorf("login fail, data:%s", string(data))
	}
	return nil
}

func (q *QBAPI) GetApplicationVersion() (string, error) {
	var version string
	err := q.getWithDecoder(apiGetAPPVersion, nil, &version, StrDec)
	if err != nil {
		return "", err
	}
	return version, nil
}

func (q *QBAPI) GetTorrentList(req *GetTorrentListRequest) (*GetTorrentListResponse, error) {
	rsp := &GetTorrentListResponse{Items: make([]*TorrentListItem, 0)}
	if err := q.getWithDecoder(apiGetTorrentList, req, &rsp.Items, JsonDec); err != nil {
		return nil, err
	}
	return rsp, nil
}
