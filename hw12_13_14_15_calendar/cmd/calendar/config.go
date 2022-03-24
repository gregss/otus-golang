package main

import (
	"flag"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level string
	File  string
}

type StorageConf struct {
	Type string
	Dsn  string
}

func NewConfig() (config Config) {
	var configFile *string
	flag.StringVar(configFile, "config", "", "verbose output")
	flag.Parse()

	yamlFile, err := ioutil.ReadFile(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return
}

// TODO
