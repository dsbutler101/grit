package ssh

import (
	"fmt"
	"io"
	"net"
	"time"
)

type streamConn struct {
	in  io.WriteCloser
	out io.ReadCloser
}

func newStreamConn(in io.WriteCloser, out io.ReadCloser) *streamConn {
	return &streamConn{
		in:  in,
		out: out,
	}
}

func (s *streamConn) Read(b []byte) (n int, err error) {
	if s.out == nil {
		return 0, nil
	}

	return s.out.Read(b)
}

func (s *streamConn) Write(b []byte) (n int, err error) {
	if s.in == nil {
		return 0, nil
	}

	return s.in.Write(b)
}

func (s *streamConn) Close() error {
	err := s.in.Close()
	if err != nil {
		return fmt.Errorf("close in stream: %w", err)
	}

	err = s.out.Close()
	if err != nil {
		return fmt.Errorf("close out stream: %w", err)
	}

	return nil
}

func (s *streamConn) LocalAddr() net.Addr {
	return nil
}

func (s *streamConn) RemoteAddr() net.Addr {
	return nil
}

func (s *streamConn) SetDeadline(t time.Time) error {
	return nil
}

func (s *streamConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (s *streamConn) SetWriteDeadline(t time.Time) error {
	return nil
}
