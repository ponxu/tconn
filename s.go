package tconn

import (
    "net"
    "sync"
)

const (
    ProtocolTCP       = "tcp"
    ProtocolWebsocket = "websocket"
)

type TConn struct {
    C        net.Conn
    O        sync.Mutex
    Buff     chan []byte
    IsClosed bool
    Tag      interface{}
}

type Handler interface {
    OnConnected(tc *TConn)
    OnClosed(tc *TConn)
    OnReceived(tc *TConn, data []byte)
}

type ProtocolAdapter interface {
    Accept(s *Server) error
    Receive(tc *TConn, v *[]byte) error
    Send(tc *TConn, v []byte) error
}

type Server struct {
    Protocol string
    Addr     string
    Adapter  ProtocolAdapter
    Handler  Handler
}

func (s *Server) Start() {
    switch s.Protocol {
    case ProtocolWebsocket:
        s.Adapter = &WSAdapter{}
        if err := s.Adapter.Accept(s); err != nil {
            panic(err)
        }
    case ProtocolTCP:
        s.Adapter = &TCPAdapter{}
        if err := s.Adapter.Accept(s); err != nil {
            panic(err)
        }
    default:
        if s.Adapter != nil {
            if err := s.Adapter.Accept(s); err != nil {
                panic(err)
            }
        } else {
            panic("Unkown Protocol: " + s.Protocol)
        }
    }
}

func (s *Server) ProcessConn(c net.Conn) {
    tc := &TConn{
        C:        c,
        Buff:     make(chan []byte, 10),
        IsClosed: false,
    }
    s.Handler.OnConnected(tc)
    go s.doWrite(tc)
    s.doRead(tc)
}

func (s *Server) doRead(tc *TConn) {
    for {
        var b []byte
        err := s.Adapter.Receive(tc, &b)
        if err != nil {
            s.doClose(tc)
            return
        }
        s.Handler.OnReceived(tc, b)
    }
}

func (s *Server) doWrite(tc *TConn) {
    for {
        b, ok := <-tc.Buff
        if !ok {
            s.doClose(tc)
            return
        }

        err := s.Adapter.Send(tc, b)
        if err != nil {
            s.doClose(tc)
            return
        }
    }
}

func (s *Server) doClose(tc *TConn) {
    if tc.IsClosed {
        return
    }

    tc.O.Lock()
    defer tc.O.Unlock()

    if !tc.IsClosed {
        tc.IsClosed = true
        tc.C.Close()   // ends read loop
        close(tc.Buff) // ends write loop
        s.Handler.OnClosed(tc)
        tc.Tag = nil
    }
}

func (tc *TConn) Send(data []byte) {
    tc.Buff <- data
}

func (tc *TConn) Close() {
    tc.C.Close()
}
