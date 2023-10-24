package rootpath

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Options .
type Options struct {
	MaxDepth uint
}

// options .
func options(opts ...func(o *Options)) *Options {
	var options = new(Options)
	for _, opt := range opts {
		opt(options)
	}
	return options
}

// Root .
type Root string

// NewRoot .
func NewRoot(root string) *Root {
	var r Root = Root(filepath.Clean(root) + string(filepath.Separator))
	return &r
}

// depth .
func (r *Root) depth(v string) uint {
	return uint(len(strings.Split(strings.TrimPrefix(v, r.String()), string(filepath.Separator))))
}

// match 检查路径是否符合要求
func (r *Root) match(v string, opts *Options) bool {
	// 为原 root 路径
	if len(v) == len(r.String()) {
		return false
	}
	if opts.MaxDepth != 0 {
		return r.depth(v) <= opts.MaxDepth
	}
	return true
}

// IsDir .
func (r *Root) IsDir() bool {
	info, err := os.Lstat(r.String())
	return err == nil && info.IsDir()
}

// String .
func (r *Root) String() string {
	return string(*r)
}

// Walk .
func (r *Root) Walk(fn func(path string, info fs.FileInfo) error, opts ...func(o *Options)) error {
	var options = options(opts...)
	return filepath.Walk(r.String(), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// depth
		if r.match(path, options) {
			// do something
			return fn(path, info)
		}
		return nil
	})
}

// Dirs .
func (r *Root) Dirs(opts ...func(o *Options)) ([]*Root, error) {
	var options = options(opts...)
	var roots = make([]*Root, 0)

	return roots, r.Walk(func(path string, info fs.FileInfo) error {
		if info.IsDir() && r.match(path, options) {
			roots = append(roots, NewRoot(path))
		}
		return nil
	})
}

// Files .
func (r *Root) Files(opts ...func(o *Options)) ([]*File, error) {
	var options = options(opts...)
	var files = make([]*File, 0)

	return files, r.Walk(func(path string, info fs.FileInfo) error {
		if !info.IsDir() && r.match(path, options) {
			files = append(files, NewFile(path))
		}
		return nil
	})
}

// File .
type File string

// String .
func (f *File) String() string {
	return string(*f)
}

// NewFile .
func NewFile(f string) *File {
	return (*File)(&f)
}
