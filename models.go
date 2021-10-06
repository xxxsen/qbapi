package qbapi

//doc: https://github.com/qbittorrent/qBittorrent/wiki/WebUI-API-(qBittorrent-4.1)

type LoginReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRsp struct {
}

type GetTorrentListRequest struct {
	Filter   *string `json:"filter,omitempty"`
	Category *string `json:"category,omitempty"`
	Tag      *string `json:"tag,omitempty"`
	Sort     *string `json:"sort,omitempty"`
	Reverse  *bool   `json:"reverse,omitempty"`
	Limit    *int    `json:"limit,omitempty"`
	Offset   *int    `json:"offset,omitempty"`
	Hashes   *string `json:"hashes,omitempty"`
}

type TorrentListItem struct {
	AddedOn           int     `json:"added_on"`
	AmountLeft        int     `json:"amount_left"`
	AutoTmm           bool    `json:"auto_tmm"`
	Availability      float64 `json:"availability"`
	Category          string  `json:"category"`
	Completed         int     `json:"completed"`
	CompletionOn      int     `json:"completion_on"`
	ContentPath       string  `json:"content_path"`
	DlLimit           int     `json:"dl_limit"`
	Dlspeed           int     `json:"dlspeed"`
	Downloaded        int     `json:"downloaded"`
	DownloadedSession int     `json:"downloaded_session"`
	Eta               int     `json:"eta"`
	FLPiecePrio       bool    `json:"f_l_piece_prio"`
	ForceStart        bool    `json:"force_start"`
	Hash              string  `json:"hash"`
	LastActivity      int     `json:"last_activity"`
	MagnetUri         string  `json:"magnet_uri"`
	MaxRatio          float64 `json:"max_ratio"`
	MaxSeedingTime    int     `json:"max_seeding_time"`
	Name              string  `json:"name"`
	NumComplete       int     `json:"num_complete"`
	NumIncomplete     int     `json:"num_incomplete"`
	NumLeechs         int     `json:"num_leechs"`
	NumSeeds          int     `json:"num_seeds"`
	Priority          int     `json:"priority"`
	Progress          float64 `json:"progress"`
	Ratio             float64 `json:"ratio"`
	RatioLimit        float64 `json:"ratio_limit"`
	SavePath          string  `json:"save_path"`
	SeedingTime       int     `json:"seeding_time"`
	SeedingTimeLimit  int     `json:"seeding_time_limit"`
	SeenComplete      int     `json:"seen_complete"`
	SeqDl             bool    `json:"seq_dl"`
	Size              int     `json:"size"`
	State             string  `json:"state"`
	SuperSeeding      bool    `json:"super_seeding"`
	Tags              string  `json:"tags"`
	TimeActive        int     `json:"time_active"`
	TotalSize         int     `json:"total_size"`
	Tracker           string  `json:"tracker"`
	UpLimit           int     `json:"up_limit"`
	Uploaded          int     `json:"uploaded"`
	UploadedSession   int     `json:"uploaded_session"`
	Upspeed           int     `json:"upspeed"`
}

type GetTorrentListResponse struct {
	Items []*TorrentListItem
}
