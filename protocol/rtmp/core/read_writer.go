package core

import (
	"bufio"
	"io"
	"sync/atomic"
)

type ReadWriter struct {
	*bufio.ReadWriter
	readError  error
	writeError error

	bytesRead  uint64
	bytesWrite uint64
}

func NewReadWriter(rw io.ReadWriter, bufSize int) *ReadWriter {
	return &ReadWriter{
		ReadWriter: bufio.NewReadWriter(bufio.NewReaderSize(rw, bufSize), bufio.NewWriterSize(rw, bufSize)),
	}
}

func (rw *ReadWriter) GetBytesCounter() (uint64, uint64) {
	r := atomic.LoadUint64(&rw.bytesRead)
	w := atomic.LoadUint64(&rw.bytesWrite)
	return r, w
}

func (rw *ReadWriter) Read(p []byte) (int, error) {
	if rw.readError != nil {
		return 0, rw.readError
	}
	n, err := io.ReadAtLeast(rw.ReadWriter, p, len(p))
	if err == nil {
		atomic.AddUint64(&rw.bytesRead, uint64(n))
	}
	rw.readError = err
	return n, err
}

func (rw *ReadWriter) ReadError() error {
	return rw.readError
}

func (rw *ReadWriter) ReadUintBE(n int) (uint32, error) {
	if rw.readError != nil {
		return 0, rw.readError
	}
	ret := uint32(0)
	for i := 0; i < n; i++ {
		b, err := rw.ReadByte()
		if err != nil {
			rw.readError = err
			return 0, err
		}
		ret = ret<<8 + uint32(b)
	}
	atomic.AddUint64(&rw.bytesRead, uint64(n))
	return ret, nil
}

func (rw *ReadWriter) ReadUintLE(n int) (uint32, error) {
	if rw.readError != nil {
		return 0, rw.readError
	}
	ret := uint32(0)
	for i := 0; i < n; i++ {
		b, err := rw.ReadByte()
		if err != nil {
			rw.readError = err
			return 0, err
		}
		ret += uint32(b) << uint32(i*8)
	}
	atomic.AddUint64(&rw.bytesRead, uint64(n))
	return ret, nil
}

func (rw *ReadWriter) Flush() error {
	if rw.writeError != nil {
		return rw.writeError
	}

	if rw.ReadWriter.Writer.Buffered() == 0 {
		return nil
	}
	return rw.ReadWriter.Flush()
}

func (rw *ReadWriter) Write(p []byte) (int, error) {
	if rw.writeError != nil {
		return 0, rw.writeError
	}
	atomic.AddUint64(&rw.bytesWrite, uint64(len(p)))
	return rw.ReadWriter.Write(p)
}

func (rw *ReadWriter) WriteError() error {
	return rw.writeError
}

func (rw *ReadWriter) WriteUintBE(v uint32, n int) error {
	if rw.writeError != nil {
		return rw.writeError
	}
	for i := 0; i < n; i++ {
		b := byte(v>>uint32((n-i-1)<<3)) & 0xff
		if err := rw.WriteByte(b); err != nil {
			rw.writeError = err
			return err
		}
	}
	atomic.AddUint64(&rw.bytesWrite, uint64(n))
	return nil
}

func (rw *ReadWriter) WriteUintLE(v uint32, n int) error {
	if rw.writeError != nil {
		return rw.writeError
	}
	for i := 0; i < n; i++ {
		b := byte(v) & 0xff
		if err := rw.WriteByte(b); err != nil {
			rw.writeError = err
			return err
		}
		v = v >> 8
	}
	atomic.AddUint64(&rw.bytesWrite, uint64(n))
	return nil
}
