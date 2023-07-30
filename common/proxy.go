package common

const (
	All = iota
	Awake
	Sleep
)

const CheckUrl = "http://%s/api/Message/WXSyncMsg"
const CheckUrlIpad = "http://%s/api/Login/HeartBeat?wxid=%s"

var Proxy []*ProxyEntry
