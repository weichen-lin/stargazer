package zapper

import (
	"net"
	"sync"

	"go.uber.org/zap/zapcore"
)

type udpLogWriter struct {
    sync.Mutex
    conn *net.UDPConn
    encoder zapcore.Encoder
}

func NewUDPLogWriter(addr string, port int, encoder zapcore.Encoder) (*udpLogWriter, error) {
    conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
        IP:   net.ParseIP(addr),
        Port: port,
    })
    if err != nil {
        return nil, err
    }

    return &udpLogWriter{
        conn:   conn,
        encoder: encoder,
    }, nil
}

func (w *udpLogWriter) Write(p []byte) (n int, err error) {
    w.Lock()
    defer w.Unlock()

    return w.conn.Write(p)
}

func (w *udpLogWriter) Close() error {
    w.Lock()
    defer w.Unlock()

    return w.conn.Close()
}