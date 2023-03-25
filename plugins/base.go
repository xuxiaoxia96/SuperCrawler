package plugins

import (
	"github.com/sirupsen/logrus"
)

type Crawler interface {
	Update()
	All()
}

var (
	factoryByName = make(map[string]func() Crawler)
)


func Register(name string, factory func() Crawler) {
	factoryByName[name] = factory
}


func Create(name string) Crawler {
	if f, ok := factoryByName[name]; ok {
		return f()
	} else {
		logrus.Error("No such func")
		return nil
	}
}
