package gzip

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

// Header .
type Header struct {
	Name    string
	Size    int64
	Mode    int64
	Modtime time.Time
}

// targzip .
type targzip struct {
	tar  *tar.Writer
	gzip *gzip.Writer

	once  sync.Once
	rwmux sync.RWMutex
	close []func() error
}

// lock .
func (tx *targzip) lock() {
	tx.rwmux.Lock()
}

// unlock .
func (tx *targzip) unlock() {
	tx.rwmux.Unlock()
}

// header return tar.Header
func (tx *targzip) hdr(opts ...func(h *Header)) (*tar.Header, error) {
	var h = &Header{Mode: int64(os.O_RDWR)}
	for _, opt := range opts {
		opt(h)
	}

	if len(h.Name) == 0 || h.Size == 0 {
		return nil, errors.New("compress failed. name or size cannot be empty.")
	}

	return &tar.Header{
		Name:    h.Name,
		Size:    h.Size,
		Mode:    h.Mode,
		ModTime: h.Modtime,
	}, nil
}

// FromReader .
func (tx *targzip) FromReader(r io.Reader, opts ...func(h *Header)) error {
	header, err := tx.hdr(opts...)
	if err != nil {
		return err
	}

	tx.lock()
	defer tx.unlock()

	// header
	if err := tx.tar.WriteHeader(header); err != nil {
		return err
	}

	// io.Copy
	if _, err := io.Copy(tx.tar, r); err != nil {
		return fmt.Errorf(`compress "%s" failed. %s`, header.Name, err.Error())
	}

	return nil
}

// WithWrite .
func (tx *targzip) WithWrite(fn func(w io.Writer) error, opts ...func(h *Header)) error {
	header, err := tx.hdr(opts...)
	if err != nil {
		return err
	}

	tx.lock()
	defer tx.unlock()

	if err := fn(tx.tar); err != nil {
		return fmt.Errorf(`compress "%s" failed. %s`, header.Name, err.Error())
	}
	return nil
}

// WithWriteAt .
func (tx *targzip) WithWriteAt(fn func(at io.WriterAt) error, opts ...func(h *Header)) error {
	header, err := tx.hdr(opts...)
	if err != nil {
		return err
	}

	tx.lock()
	defer tx.unlock()

	if err := fn(&write{tx.tar}); err != nil {
		return fmt.Errorf(`compress "%s" failed. %s`, header.Name, err.Error())
	}
	return nil
}

// Close .
func (tx *targzip) Close() error {
	tx.once.Do(func() {
		tx.tar.Close()
		tx.gzip.Close()
	})
	return nil
}

// write .
type write struct {
	io.WriteCloser
}

// WriteAt 顺序写入。不可异步写入, 否则会导致写入失败
func (w *write) WriteAt(p []byte, off int64) (int, error) {
	return w.Write(p)
}

// New .
func New(dst io.Writer) *targzip {
	tx := &targzip{gzip: gzip.NewWriter(dst)}
	tx.tar = tar.NewWriter(tx.gzip)
	return tx
}
