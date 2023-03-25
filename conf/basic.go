package conf

import (
	"os"
)

var Cfg Config

type Config struct {
	Proxy struct {
		Url string
		Token string
	}
	MongoDB struct{
		Host			string
		Port			string
		Username		string
		Password		string
		AuthSource		string
		AuthMechanism	string
	}
}


func InitConfig() {
	if os.Getenv("DEBUG") == "1"{
		// Proxy
		Cfg.Proxy.Url = "https://proxy.webshare.io/api/proxy/list/?page=1"
		Cfg.Proxy.Token = "7087845876b78a62e3cba5603902c1e92e9f847a"

		// MongoDB
		Cfg.MongoDB.Host = "localhost"
		Cfg.MongoDB.Port = "27017"
		Cfg.MongoDB.Username = ""
		Cfg.MongoDB.Password = ""
		Cfg.MongoDB.AuthSource = "admin"
		Cfg.MongoDB.AuthMechanism = "SCRAM-SHA-1"
	}

}