package ssh

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMultiReadCloser_Read(t *testing.T) {
	reader1Content := "reader 1"
	reader2Content := "reader 2"

	reader1 := io.NopCloser(bytes.NewBufferString(reader1Content))
	reader2 := io.NopCloser(bytes.NewBufferString(reader2Content))
	r := newMultiReadCloser(reader1, reader2)

	out := bytes.NewBuffer(nil)

	n, err := io.Copy(out, r)
	assert.Equal(t, int64(len(reader1Content)+len(reader2Content)), n)
	assert.NoError(t, err)

	str := out.String()
	assert.Contains(t, str, reader1Content)
	assert.Contains(t, str, reader2Content)
}

func TestMultiReadCloser_Close(t *testing.T) {
	testError1 := errors.New("test error 1")
	testError2 := errors.New("test error 2")

	tests := map[string]struct {
		mockReaders func(t *testing.T) []io.ReadCloser
		assertError func(t *testing.T, err error)
	}{
		"no errors": {
			mockReaders: func(t *testing.T) []io.ReadCloser {
				reader1 := newMockTestReadCloser(t)
				reader1.EXPECT().Close().Return(nil)
				reader2 := newMockTestReadCloser(t)
				reader2.EXPECT().Close().Return(nil)

				return []io.ReadCloser{reader1, reader2}
			},
		},
		"first reader returns error": {
			mockReaders: func(t *testing.T) []io.ReadCloser {
				reader1 := newMockTestReadCloser(t)
				reader1.EXPECT().Close().Return(assert.AnError)
				reader2 := newMockTestReadCloser(t)
				reader2.EXPECT().Close().Return(nil)

				return []io.ReadCloser{reader1, reader2}
			},
			assertError: func(t *testing.T, err error) {
				var eerr *multiReaderCloserCloseErr
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, 1, eerr.failures)
					assert.Equal(t, assert.AnError, eerr.lastErr)
				}
			},
		},
		"second reader returns error": {
			mockReaders: func(t *testing.T) []io.ReadCloser {
				reader1 := newMockTestReadCloser(t)
				reader1.EXPECT().Close().Return(nil)
				reader2 := newMockTestReadCloser(t)
				reader2.EXPECT().Close().Return(assert.AnError)

				return []io.ReadCloser{reader1, reader2}
			},
			assertError: func(t *testing.T, err error) {
				var eerr *multiReaderCloserCloseErr
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, 1, eerr.failures)
					assert.Equal(t, assert.AnError, eerr.lastErr)
				}
			},
		},
		"twi readers return error": {
			mockReaders: func(t *testing.T) []io.ReadCloser {
				reader1 := newMockTestReadCloser(t)
				reader1.EXPECT().Close().Return(testError1)
				reader2 := newMockTestReadCloser(t)
				reader2.EXPECT().Close().Return(testError2)

				return []io.ReadCloser{reader1, reader2}
			},
			assertError: func(t *testing.T, err error) {
				var eerr *multiReaderCloserCloseErr
				if assert.ErrorAs(t, err, &eerr) {
					assert.Equal(t, 2, eerr.failures)
					assert.Equal(t, testError2, eerr.lastErr)
				}
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			require.NotNil(t, tt.mockReaders, "mockReaders must be defined in test definition")

			err := newMultiReadCloser(tt.mockReaders(t)...).Close()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestStreamConn_Read(t *testing.T) {
	testContent := "test content"

	tests := map[string]struct {
		prepareReaderMock func(t *testing.T) io.ReadCloser
		assertError       func(t *testing.T, err error)
		expectedOutput    []byte
	}{
		"outputStream not provided": {
			prepareReaderMock: func(t *testing.T) io.ReadCloser {
				return nil
			},
			expectedOutput: []byte{},
		},
		"outputStream read error": {
			prepareReaderMock: func(t *testing.T) io.ReadCloser {
				m := newMockTestReadCloser(t)
				m.EXPECT().Read(mock.Anything).Return(0, assert.AnError)

				return m
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"outputStream correct read": {
			prepareReaderMock: func(t *testing.T) io.ReadCloser {
				m := newMockTestReadCloser(t)
				m.EXPECT().
					Read(mock.Anything).
					Run(func(p []byte) {
						copy(p, testContent)
					}).
					Return(len(testContent), io.EOF)

				return m
			},
			expectedOutput: []byte(testContent),
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			require.NotNil(t, tt.prepareReaderMock, "prepareReaderMock must be defined in test definition")

			sc := newStreamConn(nil, tt.prepareReaderMock(t))

			out, err := io.ReadAll(sc)
			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedOutput, out)
		})
	}
}

func TestStreamConn_Write(t *testing.T) {
	testContent := "test content"

	tests := map[string]struct {
		prepareWriterMock func(t *testing.T) io.WriteCloser
		assertError       func(t *testing.T, err error)
		expectedLen       int
	}{
		"inStream not provided": {
			prepareWriterMock: func(t *testing.T) io.WriteCloser {
				return nil
			},
		},
		"inStream write error": {
			prepareWriterMock: func(t *testing.T) io.WriteCloser {
				m := newMockTestWriteCloser(t)
				m.EXPECT().Write(mock.Anything).Return(0, assert.AnError)

				return m
			},
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
			},
		},
		"inStream correct write": {
			prepareWriterMock: func(t *testing.T) io.WriteCloser {
				m := newMockTestWriteCloser(t)
				m.EXPECT().
					Write(mock.Anything).
					Run(func(p []byte) {
						assert.Equal(t, testContent, string(p))
					}).
					Return(len(testContent), nil)

				return m
			},
			expectedLen: len(testContent),
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			require.NotNil(t, tt.prepareWriterMock, "prepareWriterMock must be defined in test definition")

			sc := newStreamConn(tt.prepareWriterMock(t), nil)
			n, err := fmt.Fprint(sc, testContent)

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expectedLen, n)
		})
	}
}

func TestStreamConn_Close(t *testing.T) {
	nilInStream := func(_ *testing.T) io.WriteCloser {
		return nil
	}

	nilOutStream := func(_ *testing.T) io.ReadCloser {
		return nil
	}

	testInStream := func(e error) func(t *testing.T) io.WriteCloser {
		return func(t *testing.T) io.WriteCloser {
			m := newMockTestWriteCloser(t)
			m.EXPECT().Close().Return(e)

			return m
		}
	}

	testOutStream := func(e error) func(t *testing.T) io.ReadCloser {
		return func(t *testing.T) io.ReadCloser {
			m := newMockTestReadCloser(t)
			m.EXPECT().Close().Return(e)

			return m
		}
	}

	tests := map[string]struct {
		inStream    func(t *testing.T) io.WriteCloser
		outStream   func(t *testing.T) io.ReadCloser
		assertError func(t *testing.T, err error)
	}{
		"no in nor out streams": {
			inStream:  nilInStream,
			outStream: nilOutStream,
		},
		"no in stream": {
			inStream:  nilInStream,
			outStream: testOutStream(nil),
		},
		"no out stream": {
			inStream:  testInStream(nil),
			outStream: nilOutStream,
		},
		"both streams provided": {
			inStream:  testInStream(nil),
			outStream: testOutStream(nil),
		},
		"inStream close error": {
			inStream:  testInStream(assert.AnError),
			outStream: nilOutStream,
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errStreamConnClosingInStream)
			},
		},
		"outStream close error": {
			inStream:  nilInStream,
			outStream: testOutStream(assert.AnError),
			assertError: func(t *testing.T, err error) {
				assert.ErrorIs(t, err, assert.AnError)
				assert.ErrorIs(t, err, errStreamConnClosingOutStream)
			},
		},
	}

	for tn, tt := range tests {
		t.Run(tn, func(t *testing.T) {
			require.NotNil(t, tt.inStream, "inStream must be defined in test definition")
			require.NotNil(t, tt.outStream, "outStream must be defined in test definition")

			sc := newStreamConn(tt.inStream(t), tt.outStream(t))

			err := sc.Close()

			if tt.assertError != nil {
				tt.assertError(t, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
