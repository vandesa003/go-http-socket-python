package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// http request struct
type HttpRequest struct {
	image string `json:"image"`
}

// http server conf.
const (
	timeout        time.Duration = 3 * time.Second
	host           string        = "0.0.0.0"
	port           int32         = 8080
	socketAddr     string        = "test.socket"
	socketConnType string        = "unix"
)

func main() {
	// init protocol.
	ptc := newProtocol(socketAddr)
	// start socket listener.
	go ptc.run()

	// fasthttp server
	server := &fasthttp.Server{
		Handler:            ptc.Handler,
		MaxRequestBodySize: 1 * 1024 * 1024 * 1024,
		DisableKeepalive:   true,
	}

	// run server with 1 goroutine.
	go func() {
		log.Info().Msg("fasthttp running")
		if err := server.ListenAndServe(fmt.Sprintf("%s:%d", host, port)); err != nil {
			log.Fatal().Err(err).Msg("fasthttp error")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	log.Info().Msg("get the terminate signal, try to shutdown service")
	shutdown := make(chan bool, 1)
	go func(exit chan bool) {
		if err := server.Shutdown(); err != nil {
			log.Fatal().Err(err).Msg("error in fasthttp shutdown")
		}
		exit <- true
	}(shutdown)

	select {
	case <-shutdown:
		log.Info().Msg("fasthttp shutdown successful")
	case <-time.After(timeout):
		log.Warn().Msg("cannot shutdown fasthttp")
	}
}
