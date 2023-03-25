package common

import (
	"SuperCrawler/conf"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"github.com/sirupsen/logrus"
	"net/url"
	"time"
)

type Ports struct {
	Http 		int
	Socks5 		int
}

type ProxyItem struct{
	Username 				string
	Password     			string
	ProxyAddress 			string	`json:"proxy_address"`
	Ports         			Ports
	Valid 					bool
	LastVerification 		string 	`json:"last_verification"`
	CountryCode 			string
	CountryCodeConfidence 	float64 `json:"country_code_confidence"`
	CityName 				string 	`json:"city_name"`
}

type ProxyRes struct{
	Count 		int
	Next 		interface{}
	Previous 	interface{}
	Results 	[]ProxyItem
}

func GetAProxy() *url.URL{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	retried := 0
	// retry 3 times
	for retried < 3{
		if retried > 0{
			logrus.Warning(fmt.Sprintf("[GetAProxy] Retry times: %d", retried))
			time.Sleep(time.Duration(3)*time.Second)
		}
		// fetch api
		req, err := http.NewRequest("GET", conf.Cfg.Proxy.Url, nil)
		if err != nil {
			logrus.Error(fmt.Sprintf("[GetAProxy fetch] %s", err))
			retried += 1
			continue
		}
		req.Header.Set("Authorization", conf.Cfg.Proxy.Token)
		resp, err := client.Do(req)
		if err != nil {
			logrus.Error(fmt.Sprintf("[GetAProxy do] %s", err))
			retried += 1
			continue
		}
		defer resp.Body.Close()

		// parse result
		respByte, _ := ioutil.ReadAll(resp.Body)
		proxyRes := ProxyRes{}
		err = json.Unmarshal(respByte, &proxyRes)
		if err != nil {
			logrus.Error(fmt.Sprintf("[GetAProxy] %s", err))
			retried += 1
			continue
		}
		if proxyRes.Count > 0 {
			rand.Seed(time.Now().Unix())
			selected_proxy := proxyRes.Results[rand.Intn(proxyRes.Count-1)]
			proxy, _ := url.Parse(fmt.Sprintf("http://%s:%s@%s:%d", selected_proxy.Username, selected_proxy.Password, selected_proxy.ProxyAddress, selected_proxy.Ports.Http))
			return proxy
		}else{
			continue
		}
	}

	return nil
}