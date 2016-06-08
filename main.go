package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/mikespook/golib/log"
	"github.com/mikespook/golib/signal"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "Config file")
	flag.Parse()
}

func main() {
	c, err := InitConfig(configFile)
	if err != nil {
		log.Error(err)
		return
	}

	http.HandleFunc("/", qrGen)

	signal.Bind(os.Interrupt, func() uint {
		log.Message("Exit")
		return signal.BreakExit
	})

	log.Messagef("Listening: %s", c.Addr)
	go func() {
		if err := http.ListenAndServe(c.Addr, nil); err != nil {
			log.Error(err)
			if err := signal.Send(os.Getpid(), os.Interrupt); err != nil {
				panic(err)
			}
			return
		}
	}()
	signal.Wait()
}
