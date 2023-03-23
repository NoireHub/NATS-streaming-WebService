package main

import (
	"flag"
	"log"
	"os/signal"
	"os"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/NoireHub/NATS-streaming-WebService/internal/app/webservice"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath,"","configs/webservice.toml","path to config file")
}

func handleCrash() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Fatal("interrupt")
		os.Exit(0)
	}()
}

func main() {
	flag.Parse()
	handleCrash()

	config:= webservice.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	} 

	if err := webservice.Start(config); err != nil {
		log.Fatal(err)
	}
}