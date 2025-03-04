package main

import (
	"flag"
	"fmt"
	"github.com/Xhofe/alist/bootstrap"
	"github.com/Xhofe/alist/conf"
	"github.com/Xhofe/alist/model"
	"github.com/Xhofe/alist/server"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func init() {
	flag.StringVar(&conf.ConfigFile, "conf", "data/config.json", "config file")
	flag.BoolVar(&conf.Debug, "debug", false, "start with debug mode")
	flag.BoolVar(&conf.Version, "version", false, "print version info")
	flag.BoolVar(&conf.Password, "password", false, "print current password")
	flag.Parse()
}

func Init() bool {
	bootstrap.InitLog()
	bootstrap.InitConf()
	bootstrap.InitCron()
	bootstrap.InitModel()
	if conf.Password {
		pass, err := model.GetSettingByKey("password")
		if err != nil {
			log.Errorf(err.Error())
			return false
		}
		log.Infof("current password: %s", pass.Value)
		return false
	}
	bootstrap.InitSettings()
	bootstrap.InitAccounts()
	bootstrap.InitCache()
	return true
}

func main() {
	if conf.Version {
		fmt.Printf("Built At: %s\nGo Version: %s\nAuthor: %s\nCommit ID: %s\nVersion: %s\n", conf.BuiltAt, conf.GoVersion, conf.GitAuthor, conf.GitCommit, conf.GitTag)
		return
	}
	if !Init() {
		return
	}
	if !conf.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	server.InitApiRouter(r)
	base := fmt.Sprintf("%s:%d", conf.Conf.Address, conf.Conf.Port)
	log.Infof("start server @ %s", base)
	var err error
	if conf.Conf.Https {
		err = r.RunTLS(base, conf.Conf.CertFile, conf.Conf.KeyFile)
	} else {
		err = r.Run(base)
	}
	if err != nil {
		log.Errorf("failed to start: %s", err.Error())
	}
}
