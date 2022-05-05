
package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"go.uber.org/config"
	"ascendex.io/act-aws-lambda-s3/apiserver"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "./apiserver.yaml", "config file path")
}

func initGroupCache() {

}

func main() {

	a := apiserver.App{}

	var log = logrus.New()
	opt := config.File(confPath)
	yaml, err := config.NewYAML(opt)
	if err != nil {
		log.Error(err)
	}

	cfg := &apiserver.Config{}
	if err := yaml.Get("apiserver").Populate(cfg); err != nil {
		log.Error(err)
	}
	if err = a.Initialize(cfg); err != nil {
		log.Error(err)
	}

	log.Info("server start")
	a.Run()
}
