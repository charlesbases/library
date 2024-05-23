package s3

import (
	"math"
	"testing"
	"time"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library"
	"github.com/charlesbases/library/storage"
)

var (
	c storage.Client

	clientBOS = func() (storage.Client, error) {
		endpoint := "s3.bcebos.cncq.icpc.changan.com"
		accessKey := "437e8bdc81b14da796789da67667dd52"
		secretKey := "9eb1e112d8a144a8ab125020cf6e7403"
		return NewClient(endpoint, accessKey, secretKey)
	}

	clientMinIO = func() (storage.Client, error) {
		endpoint := "10.64.21.64:32607"
		accessKey := "AINhzLWnmdnD70Ve"
		secretKey := "zxcfd6I3UeoaXrMHigA48YEbHy39Hmji"
		return NewClient(endpoint, accessKey, secretKey)
	}

	keyPrefix = func(v string) string {
		return "testdata/data/" + v
	}

	mxdata = "mxdata"
)

const (
	// 上传文件夹本地路径
	root = `C:\Users\sun\Desktop\test\auth`
	// 远程文件下载带本地
	folder = `C:\Users\sun\Desktop\test\copy`
	// s3路径
	src = "auth/"
	// 对象拷贝时的路径
	dst = "copy/"
)

func init() {
	client, err := clientBOS()
	if err != nil {
		logger.Fatal(err)
	}
	c = client
}

// PutObject
// GetObject
// DelObject
// IsExist
func TestS3(t *testing.T) {
	t.Run("PutObject", func(t *testing.T) {
		key := "int"
		t.Run(key, func(t *testing.T) {
			stderr(c.PutObject(storage.InputNumber(mxdata, keyPrefix(key), time.Now().UnixMicro())))
		})

		key = "bool"
		t.Run(key, func(t *testing.T) {
			stderr(c.PutObject(storage.InputBoolean(mxdata, keyPrefix(key), true)))
		})

		key = "float"
		t.Run(key, func(t *testing.T) {
			stderr(c.PutObject(storage.InputNumber(mxdata, keyPrefix(key), math.Pi)))
		})

		key = "string"
		t.Run(key, func(t *testing.T) {
			stderr(c.PutObject(storage.InputString(mxdata, keyPrefix(key), library.NowString())))
		})
	})

	t.Run("GetObject", func(t *testing.T) {
		key := "int"
		t.Run(key, func(t *testing.T) {
			if outhook, err := c.GetObject(mxdata, keyPrefix(key)); err != nil {
				stderr(err)
			} else {
				var v interface{}
				if err := outhook.Fetch(func(output storage.ObjectOutput) error {
					return output.Decode(&v)
				}); err != nil {
					stderr(err)
				} else {
					logger.Debugf("[aws-s3]: %v", v)
				}
			}
		})

		key = "bool"
		t.Run(key, func(t *testing.T) {
			if outhook, err := c.GetObject(mxdata, keyPrefix(key)); err != nil {
				stderr(err)
			} else {
				var v bool
				if err := outhook.Fetch(func(output storage.ObjectOutput) error {
					return output.Decode(&v)
				}); err != nil {
					stderr(err)
				} else {
					logger.Debugf("[aws-s3]: %v", v)
				}
			}
		})

		key = "float"
		t.Run(key, func(t *testing.T) {
			if outhook, err := c.GetObject(mxdata, keyPrefix(key)); err != nil {
				stderr(err)
			} else {
				var v float64
				if err := outhook.Fetch(func(output storage.ObjectOutput) error {
					return output.Decode(&v)
				}); err != nil {
					stderr(err)
				} else {
					logger.Debugf("[aws-s3]: %v", v)
				}
			}
		})

		key = "string"
		t.Run(key, func(t *testing.T) {
			if outhook, err := c.GetObject(mxdata, keyPrefix(key)); err != nil {
				stderr(err)
			} else {
				var v string
				if err := outhook.Fetch(func(output storage.ObjectOutput) error {
					return output.Decode(&v)
				}); err != nil {
					stderr(err)
				} else {
					logger.Debugf("[aws-s3]: %v", v)
				}
			}
		})

		key = "notfound"
		t.Run(key, func(t *testing.T) {
			if outhook, err := c.GetObject(mxdata, keyPrefix(key)); err != nil {
				stderr(err)
			} else {
				var v string
				if err := outhook.Fetch(func(output storage.ObjectOutput) error {
					return output.Decode(&v)
				}); err != nil {
					stderr(err)
				} else {
					logger.Debugf("[aws-s3]: %v", v)
				}
			}
		})
	})

	t.Run("DelObject", func(t *testing.T) {
		key := "int"
		t.Run(key, func(t *testing.T) {
			stderr(c.DelObject(mxdata, keyPrefix(key)))
		})

		key = "bool"
		t.Run(key, func(t *testing.T) {
			stderr(c.DelObject(mxdata, keyPrefix(key)))
		})

		key = "float"
		t.Run(key, func(t *testing.T) {
			stderr(c.DelObject(mxdata, keyPrefix(key)))
		})

		key = "string"
		t.Run(key, func(t *testing.T) {
			stderr(c.DelObject(mxdata, keyPrefix(key)))
		})

		key = "notfound"
		t.Run(key, func(t *testing.T) {
			stderr(c.DelObject(mxdata, keyPrefix(key)))
		})
	})
}

// PutFolder
// Copy
// GetObjectsWithIterator
// Downloads
// DelObjectsWithPrefix
// IsExist
// Presign
func BenchmarkS3(b *testing.B) {
	var bench = func(f func()) {
		b.ResetTimer()
		f()
		b.StopTimer()
	}

	b.Run("PutFolder", func(b *testing.B) {
		bench(func() {
			stderr(c.PutFolder(mxdata, src, root))
		})
	})

	b.Run("Copy", func(b *testing.B) {
		stderr(c.Copy(storage.PositionRemote(mxdata, src), storage.PositionRemote(mxdata, dst)))
	})

	b.Run("GetObjectsWithIterator", func(b *testing.B) {
		bench(func() {
			var total int
			stderr(c.GetObjectsWithIterator(mxdata, keyPrefix(""), func(keys []*string) error {
				total += len(keys)
				return nil
			}))
			logger.Debug(total)
		})
	})

	b.Run("Downloads", func(b *testing.B) {
		bench(func() {
			stderr(c.Downloads(mxdata, dst, folder))
		})
	})

	b.Run("DelObjectsWithPrefix", func(b *testing.B) {
		bench(func() {
			stderr(c.DelObjectsWithPrefix(mxdata, dst))
		})
	})

	b.Run("IsExist", func(b *testing.B) {
		bench(func() {
			exists, err := c.IsExist(mxdata, dst)
			stderr(err)
			logger.Debug(exists)
		})
	})

	b.Run("Presign", func(b *testing.B) {
		bench(func() {
			url, err := c.Presign(mxdata, dst)
			stderr(err)
			logger.Debug(url)
		})
	})
}

// stderr .
func stderr(err error) {
	if err != nil {
		logger.Error(err)
	}
}
