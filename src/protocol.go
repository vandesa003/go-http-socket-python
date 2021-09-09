package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

// socket protocol
type Protocol struct {
	address  string
	listener net.Listener
	task     chan []byte
}

// init protocol.
func newProtocol(addr string) *Protocol {
	if _, err := os.Stat(addr); err == nil {
		log.Info().Str("addr", addr).Msg("socket already exists, try to delete it")
		if err := os.Remove(addr); err != nil {
			log.Fatal().Err(err).Msg("remove socket file error")
		}
	}
	listener, err := net.Listen("unix", addr)
	if err != nil {
		log.Fatal().Err(err).Str("addr", addr).Msg("cannot build socket connection.")
	}
	log.Info().Str("addr", addr).Msg("listen to socket")
	return &Protocol{
		address:  addr,
		listener: listener,
		task:     make(chan []byte, 10),
	}
}

// run socket server.
func (p *Protocol) run() {
	for {
		conn, err := p.listener.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "use of closed network connection") {
				log.Info().Msg("close the listener")
				break
			}
			log.Error().Err(err).Msg("Protocol listener accept error")
		}
		log.Info().
			Str("local", conn.LocalAddr().String()).
			Str("remote", conn.RemoteAddr().String()).
			Msg("accept connection")

		go p.communicate(conn)
	}
}

// socket communication.
func (p *Protocol) communicate(conn net.Conn) {
	for {
		select {
		case d := <-p.task:
			conn.Write(d)
			fmt.Println("sent task to socket.")
		}
	}
}

// fasthttp handler.
func (p *Protocol) Handler(ctx *fasthttp.RequestCtx) {
	var req HttpRequest
	body := ctx.PostBody()
	p.task <- body
	if err := json.Unmarshal(body, &req); err != nil {
		ctx.SetContentType("application/json")
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("internal error."))
	}
	fmt.Println(req.image)
	ctx.SetBody([]byte("hello world!"))
}
