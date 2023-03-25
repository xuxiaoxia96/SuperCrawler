package main

import (
	"SuperCrawler/conf"
	"SuperCrawler/plugins"
	"SuperCrawler/vars"
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

func main(){
	// flag parse
	flag.Parse()
	if *vars.Version {
		fmt.Println(vars.VersionInfo)
		return
	}

	// init config
	conf.InitConfig()

	var registers []string
	// register list
	if len(*vars.Target) == 0{
		registers = []string{"Nation"}
	}else{
		registers = strings.Split(*vars.Target, ",")
	}

	for _,register := range registers{
		module := plugins.Create(register)
		if *vars.Mode == "update"{
			module.Update()
		}else if *vars.Mode == "all"{
			module.All()
		}else{
			logrus.Error("[Main][Module Func] No such registered func")
		}
	}

}
