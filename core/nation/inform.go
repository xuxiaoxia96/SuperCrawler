package nation

import (
	"SuperCrawler/common"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type InformDucoment struct {
	UniqId				string				`bson:"uniq_id"`
	District 		string				`bson:"district"`
	SubDistrict 	string				`bson:"sub_district"`
	ExamType 		string				`bson:"exam_type"`
	InfoType		string				`bson:"info_type"`
	Title 			string				`bson:"title"`
	Content 		string				`bson:"content"`
	InsertTime 		time.Time			`bson:"insert_time"`
	UpdateTime 		time.Time 			`bson:"update_time"`
}

type InformRespObj struct {
	Articlelist 		[]InformArticleObj
	CacheInfo			interface{}
	TotalCount			int64
}

type InformArticleObj struct {
	Id					string
	CmsArticleColumnId	string
	ArticleTitle		string
	Pstrtime 			int64
	Fj					string
	ArticleType			string
	ArticleUrl			string
	IsFirst				string
	IsIntab				string
	Ahd011 				string
}

type InformArticleContObj struct {
	ArticleTitle		string
	Ctime				int64
	Content				string
	Ahd011				interface{}
}

type InformContentObj struct {
	CacheInfo			interface{}
	Article				InformArticleContObj
	ResourceList		interface{}
}



func (informDucoment *InformDucoment)getHeaders() map[string]string {
	headers := map[string]string{
		"User-Agent": common.RandUA(),
	}

	return headers
}


func (informDucoment *InformDucoment)getInformBase(uri string, proxy *url.URL) *InformRespObj{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout: 10 * time.Second,
		Proxy:           http.ProxyURL(proxy),
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
		return nil
	}
	headers := informDucoment.getHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
		return nil
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	nationRespObject := InformRespObj{}
	err = json.Unmarshal(respByte, &nationRespObject)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
		return nil
	}

	return &nationRespObject
}


func (informDucoment *InformDucoment)getInformTotalCnt(uri string, proxy *url.URL) int64{
	var totalCnt int64
	nationRespObject := informDucoment.getInformBase(uri, proxy)
	if nationRespObject != nil{
		totalCnt = nationRespObject.TotalCount
	}else{
		totalCnt = 0
	}

	return totalCnt
}


func (informDucoment *InformDucoment)getPageInformSingle(uri string, proxy *url.URL) ([]interface{}){
	var result []interface{}

	nationRespObject := informDucoment.getInformBase(uri, proxy)
	if nationRespObject == nil{
		return result
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout: 10 * time.Second,
		Proxy:           http.ProxyURL(proxy),
	}
	client := &http.Client{Transport: tr}
	for _, article := range nationRespObject.Articlelist{
		nationDocument := new(InformDucoment)
		nationDocument.UniqId = article.Id
		nationDocument.District = "国家"
		nationDocument.SubDistrict = "国家"
		nationDocument.InfoType = "通知公示"
		nationDocument.ExamType = common.AddExamTag(article.ArticleTitle)
		nationDocument.Title = article.ArticleTitle
		articleUri := fmt.Sprintf("http://dl.scs.gov.cn/api/article/%s", article.Id)
		articleReq, err := http.NewRequest("GET", articleUri, nil)
		if err != nil{
			logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
			continue
		}
		headers := informDucoment.getHeaders()
		for k,v := range headers {
			articleReq.Header.Set(k, v)
		}
		articleResp, err := client.Do(articleReq)
		if err != nil {
			logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
			continue
		}
		defer articleResp.Body.Close()
		articleRespByte, _ := ioutil.ReadAll(articleResp.Body)
		nationContentObject := InformContentObj{}
		err = json.Unmarshal(articleRespByte, &nationContentObject)
		if err != nil{
			logrus.Error(fmt.Sprintf("[Nation][getPageInformSingle] %s", err))
			continue
		}
		nationDocument.Content = nationContentObject.Article.Content

		result = append(result, *nationDocument)
		//time.Sleep(time.Duration(10)*time.Second)

		break
	}

	return result
}

func GetPageInformAll() []interface{}{
	informDucoment := InformDucoment{}
	var result []interface{}
	baseUri := "http://dl.scs.gov.cn/api/article/articlelist/all/8a81f3247b82076f017b95cb49ad002e/0000000062b7b2b60162bccf480c000a/%d"
	proxy := common.GetAProxy()
	totalCnt := informDucoment.getInformTotalCnt(fmt.Sprintf(baseUri, 1), proxy)
	EachPageCnt := 30
	pages := int(totalCnt) / EachPageCnt
	logrus.Info(fmt.Sprintf("[Nation][All] Total Pages: %d", pages))
	for i:=1;i<pages+1;i++{
		proxy := common.GetAProxy()
		singleResult := informDucoment.getPageInformSingle(fmt.Sprintf(baseUri, i), proxy)
		result = append(result, singleResult...)
	}

	return result
}

func GetPageInformUpdate() []interface{}{
	informDucoment := InformDucoment{}
	var result []interface{}
	baseUri := "http://dl.scs.gov.cn/api/article/articlelist/all/8a81f3247b82076f017b95cb49ad002e/0000000062b7b2b60162bccf480c000a/%d"
	for i:=1;i<2;i++{
		proxy := common.GetAProxy()
		singleResult := informDucoment.getPageInformSingle(fmt.Sprintf(baseUri, i), proxy)
		result = append(result, singleResult...)
	}

	return result
}
