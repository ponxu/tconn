package tconn

import (
    "code.google.com/p/go.net/websocket"
    "net/http"
    _ "net/http/pprof"
)

type WSAdapter struct{}

func (a *WSAdapter) Accept(s *Server) error {
    http.Handle("/", websocket.Handler(func(c *websocket.Conn) {
        s.ProcessConn(c)
    }))
    return http.ListenAndServe(s.Addr, nil)
}

func (a *WSAdapter) Receive(tc *TConn, v *[]byte) error {
    return websocket.Message.Receive(tc.C.(*websocket.Conn), v)
}

func (a *WSAdapter) Send(tc *TConn, v []byte) error {
    return websocket.Message.Send(tc.C.(*websocket.Conn), v)
}
