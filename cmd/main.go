package main

import (
	"flag"
	"fmt"
	"github.com/axliupore/judge/config"
	"github.com/axliupore/judge/internal/http"
	"github.com/axliupore/judge/internal/nsq"
	"github.com/axliupore/judge/internal/ws"
	"github.com/axliupore/judge/pkg/cache"
	"github.com/axliupore/judge/pkg/log"
	shttp "net/http"
	_ "net/http/pprof"
	"os"
)

func main() {

	log.InitLogger()

	var h int
	var w int
	var n string

	flag.IntVar(&h, "http", 0, "http port")
	flag.IntVar(&w, "ws", 0, "websocket port")
	flag.StringVar(&n, "nsqd", "", "nsqd address")
	flag.Parse()

	var param bool
	flag.Visit(func(f *flag.Flag) {
		param = true
	})
	if !param {
		printUsage()
	}

	conf := &config.Config{
		Http: config.Http{
			Port: h,
		},
		Ws: config.Ws{
			Port: w,
		},
		Nsq: config.Nsq{
			Address: n,
		},
	}

	c, err := cache.New()
	if err != nil {
		panic(err)
	}

	app, cleanup, err := InitializeApp(conf, c)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	go func() {
		log.Logger.Info("Starting pprof server on :6060")
		if err = shttp.ListenAndServe(":6060", nil); err != nil {
			log.Logger.Error("Failed to start pprof server", err)
		}
	}()

	app.Run()
}

type App struct {
	Http *http.Server
	Ws   *ws.Server
	Nsq  *nsq.Server
}

func (a *App) Run() {

	if a.Http != nil {
		go a.Http.Start()
	}
	if a.Ws != nil {

	}
	if a.Nsq != nil {

	}

	select {}
}

func printUsage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <args>\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}
