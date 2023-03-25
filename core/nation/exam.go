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

type ExamDucoment struct {
	UniqId			string				`bson:"uniq_id"`
	District 		string				`bson:"district"`
	SubDistrict 	string				`bson:"sub_district"`
	ExamType 		string				`bson:"exam_type"`
	InfoType		string				`bson:"info_type"`
	Title 			string				`bson:"title"`
	Content 		string				`bson:"content"`
	InsertTime 		time.Time			`bson:"insert_time"`
	UpdateTime 		time.Time 			`bson:"update_time"`
}

type ExamRespObj struct {
	Articlelist 		[]ExamArticleObj
	CacheInfo			interface{}
	TotalCount			int64
}

type ExamArticleObj struct {
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

type ExamArticleContObj struct {
	ArticleTitle		string
	Ctime				int64
	Content				string
	Ahd011				interface{}
}

type ExamContentObj struct {
	CacheInfo			interface{}
	Article				ExamArticleContObj
	ResourceList		interface{}
}



func (examDucoment *ExamDucoment)getHeaders() map[string]string {
	headers := map[string]string{
		"User-Agent": common.RandUA(),
	}

	return headers
}


func (examDucoment *ExamDucoment)getExamBase(uri string, proxy *url.URL) *ExamRespObj{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout: 10 * time.Second,
		Proxy:           http.ProxyURL(proxy),
	}
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
		return nil
	}
	headers := examDucoment.getHeaders()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
		return nil
	}
	defer resp.Body.Close()
	respByte, _ := ioutil.ReadAll(resp.Body)
	nationRespObject := ExamRespObj{}
	err = json.Unmarshal(respByte, &nationRespObject)
	if err != nil {
		logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
		return nil
	}

	return &nationRespObject
}


func (examDucoment * ExamDucoment)getExamTotalCnt(uri string, proxy *url.URL) int64{
	var totalCnt int64
	nationRespObject := examDucoment.getExamBase(uri, proxy)
	if nationRespObject != nil{
		totalCnt = nationRespObject.TotalCount
	}else{
		totalCnt = 0
	}

	return totalCnt
}


func (examDucoment *ExamDucoment)getPageExamSingle(uri string, proxy *url.URL) ([]interface{}){
	var result []interface{}

	nationRespObject := examDucoment.getExamBase(uri, proxy)
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
		nationDocument := new(ExamDucoment)
		nationDocument.UniqId = article.Id
		nationDocument.District = "国家"
		nationDocument.SubDistrict = "国家"
		nationDocument.InfoType = "招考公告"
		nationDocument.ExamType = common.AddExamTag(article.ArticleTitle)
		nationDocument.Title = article.ArticleTitle
		articleUri := fmt.Sprintf("http://dl.scs.gov.cn/api/article/%s", article.Id)
		articleReq, err := http.NewRequest("GET", articleUri, nil)
		if err != nil{
			logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
			continue
		}
		headers := examDucoment.getHeaders()
		for k,v := range headers {
			articleReq.Header.Set(k, v)
		}
		articleResp, err := client.Do(articleReq)
		if err != nil {
			logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
			continue
		}
		defer articleResp.Body.Close()
		articleRespByte, _ := ioutil.ReadAll(articleResp.Body)
		nationContentObject := ExamContentObj{}
		err = json.Unmarshal(articleRespByte, &nationContentObject)
		if err != nil{
			logrus.Error(fmt.Sprintf("[Nation][getPageExamSingle] %s", err))
			continue
		}
		nationDocument.Content = nationContentObject.Article.Content

		result = append(result, *nationDocument)
		//time.Sleep(time.Duration(10)*time.Second)

		break
	}

	return result
}

func GetPageExamAll() []interface{}{
	examDucoment := ExamDucoment{}
	var result []interface{}
	baseUri := "http://dl.scs.gov.cn/api/article/articlelist/all/8a81f3247b82076f017b95cb49ad002e/0000000062b7b2b60162bccdd5860007/%d"
	proxy := common.GetAProxy()
	totalCnt := examDucoment.getExamTotalCnt(fmt.Sprintf(baseUri, 1), proxy)
	EachPageCnt := 30
	pages := int(totalCnt) / EachPageCnt + 1
	logrus.Info(fmt.Sprintf("[Nation][All] Total Pages: %d", pages))
	for i:=1;i<pages+1;i++{
		proxy := common.GetAProxy()
		singleResult := examDucoment.getPageExamSingle(fmt.Sprintf(baseUri, i), proxy)
		result = append(result, singleResult...)
	}

	return result
}

func GetPageExamUpdate() []interface{}{
	examDucoment := ExamDucoment{}
	var result []interface{}
	baseUri := "http://dl.scs.gov.cn/api/article/articlelist/all/8a81f3247b82076f017b95cb49ad002e/0000000062b7b2b60162bccdd5860007/%d"
	for i:=1;i<2;i++{
		proxy := common.GetAProxy()
		singleResult := examDucoment.getPageExamSingle(fmt.Sprintf(baseUri, i), proxy)
		result = append(result, singleResult...)
	}

	return result
}
