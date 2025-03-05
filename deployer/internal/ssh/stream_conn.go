package ssh

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

var (
	errStreamConnClosingInStream  = errors.New("closing inStream")
	errStreamConnClosingOutStream = errors.New("closing outStream")
)

//go:generate mockery --name=testReadCloser --inpackage --with-expecter

// testReadCloser is a dummy interface definition for generating a mock
// for io.ReadCloser for unit tests
type testReadCloser interface {
	io.ReadCloser
}

//go:generate mockery --name=testWriteCloser --inpackage --with-expecter

// testReadCloser is a dummy interface definition for generating a mock
// for io.WriteCloser for unit tests
type testWriteCloser interface {
	io.WriteCloser
}

type multiReaderCloserCloseErr struct {
	failures int
	lastErr  error
}

func (e *multiReaderCloserCloseErr) Error() string {
	return fmt.Sprintf("encountered %d error(s) when closing multi reader sources; last was: %v", e.failures, e.lastErr)
}

type multiReadCloser struct {
	sources []io.ReadCloser
	r       io.Reader
}

func newMultiReadCloser(readers ...io.ReadCloser) io.ReadCloser {
	r := make([]io.Reader, len(readers))
	for no, in := range readers {
		r[no] = in
	}

	return &multiReadCloser{
		r:       io.MultiReader(r...),
		sources: readers,
	}
}

func (m *multiReadCloser) Read(p []byte) (n int, err error) {
	return m.r.Read(p)
}

func (m *multiReadCloser) Close() error {
	var lastErr error

	failures := 0
	for _, source := range m.sources {
		err := source.Close()
		if err != nil {
			lastErr = err
			failures++
		}
	}

	if lastErr != nil {
		return &multiReaderCloserCloseErr{
			failures: failures,
			lastErr:  lastErr,
		}
	}

	return nil
}

type streamConn struct {
	inStream  io.WriteCloser
	outStream io.ReadCloser
}

func newStreamConn(inStream io.WriteCloser, outStream io.ReadCloser) *streamConn {
	return &streamConn{
		inStream:  inStream,
		outStream: outStream,
	}
}

func (s *streamConn) Read(b []byte) (n int, err error) {
	if s.outStream == nil {
		return 0, io.EOF
	}

	return s.outStream.Read(b)
}

func (s *streamConn) Write(b []byte) (n int, err error) {
	if s.inStream == nil {
		return 0, nil
	}

	return s.inStream.Write(b)
}

func (s *streamConn) Close() error {
	if s.inStream != nil {
		err := s.inStream.Close()
		if err != nil {
			return fmt.Errorf("%w: %w", errStreamConnClosingInStream, err)
		}
	}

	if s.outStream != nil {
		err := s.outStream.Close()
		if err != nil {
			return fmt.Errorf("%w: %w", errStreamConnClosingOutStream, err)
		}
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
