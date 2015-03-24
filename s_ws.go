package tconn

import (
    "code.google.com/p/go.net/websocket"
    "log"
    "net/http"
    _ "net/http/pprof"
)

type WSAdapter struct{}

func (a *WSAdapter) Accept(s *Server) error {
    if s.WSEnterPoint == "" {
        s.WSEnterPoint = "/"
    }
    http.Handle(s.WSEnterPoint, websocket.Handler(func(c *websocket.Conn) {
        s.ProcessConn(c)
    }))
    log.Println("WebSocket accept at：", s.WSEnterPoint, s.Addr)
    return http.ListenAndServe(s.Addr, nil)
}

func (a *WSAdapter) Receive(tc *TConn, v *[]byte) error {
    return websocket.Message.Receive(tc.C.(*websocket.Conn), v)
}

func (a *WSAdapter) Send(tc *TConn, v []byte) error {
    return websocket.Message.Send(tc.C.(*websocket.Conn), v)
}
