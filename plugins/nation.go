package plugins

import (
	"SuperCrawler/common"
	"SuperCrawler/core/nation"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NationPlugin struct {

}


func (nationPlugin *NationPlugin) Update(){
	var items []interface{}
	// inform
	items = nation.GetPageInformUpdate()
	nationPlugin.saveDB(items)
	// exam
	items = nation.GetPageExamUpdate()
	nationPlugin.saveDB(items)
}


func (nationPlugin *NationPlugin) All(){
	var items []interface{}
	// inform
	//items = nation.GetPageInformAll()
	//nationPlugin.saveDB(items)
	// exam
	items = nation.GetPageExamAll()
	nationPlugin.saveDB(items)
}


func (nationPlugin *NationPlugin)saveDB(items []interface{}){
	if len(items) == 0{
		logrus.Warning("[Nation][saveToDB] Crawl 0 items")
		return
	}
	cli := common.GetMgoCli()
	defer func() {
		if err := cli.Disconnect(context.TODO()); err != nil {
			logrus.Error(fmt.Sprintf("[Nation][CloseMongoCli] %s", err))
		}
	}()
	coll := cli.Database("ashore").Collection("nation")
	// No order insert, enhancement
	opts := options.InsertMany().SetOrdered(false)
	result, _ := coll.InsertMany(context.TODO(), items, opts)

	var insertCnt int
	if result != nil{
		insertCnt = len(result.InsertedIDs)
	}else{
		insertCnt = 0
	}
	logrus.Info(fmt.Sprintf("[Nation][saveToDB] Crawl & Insert %d items", insertCnt))
}

func init() {
	Register("Nation", func() Crawler {
		return new(NationPlugin)
	})
}

