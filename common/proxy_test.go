package common

import (
	"SuperCrawler/conf"
	"testing"
)

func TestGetAProxy(t *testing.T) {
	conf.InitConfig()
	proxy := GetAProxy()
	t.Logf("Get Proxy: %s", proxy)
}
