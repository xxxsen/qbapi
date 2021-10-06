package qbapi

const (
	//login
	apiLogin = "/api/v2/auth/login"
	//
	apiGetAPPVersion      = "/api/v2/app/version"
	apiGetAPIVersion      = "/api/v2/app/webapiVersion"
	apiGetBuildInfo       = "/api/v2/app/buildInfo"
	apiShutdownAPP        = "/api/v2/app/shutdown"
	apiGetAPPPerf         = "/api/v2/app/preferences"
	apiSetAPPPref         = "/api/v2/app/setPreferences"
	apiGetDefaultSavePath = "/api/v2/app/defaultSavePath"
	//Log
	apiGetLog     = "/api/v2/log/main"
	apiGetPeerLog = "/api/v2/log/peers"
	//sync
	apiGetMainData        = "/api/v2/sync/maindata"
	apiGetTorrentPeerData = "/api/v2/sync/torrentPeers"
	//transfer info
	apiGetGlobalTransferInfo  = "/api/v2/transfer/info"
	apiGetAltSpeedLimitState  = "/api/v2/transfer/speedLimitsMode"
	apiToggleAltSpeedLimits   = "/api/v2/transfer/toggleSpeedLimitsMode"
	apiGetGlobalDownloadLimit = "/api/v2/transfer/downloadLimit"

	//
	apiGetTorrentList = "/api/v2/torrents/info"
)
