package storage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/pkg/errors"

	"google.golang.org/protobuf/proto"

	"github.com/charlesbases/library/content"
)

// ErrNoSuchKey .
var ErrNoSuchKey = errors.New("NoSuchKey: The specified key does not exist.")

// Client .
type Client interface {
	// PutObject 上传单个对象
	PutObject(input ObjectInput, opts ...func(o *PutOptions)) error
	// PutFolder 上传本地文件夹，存储路径为 'prefix/root/xxx'
	PutFolder(bucket string, prefix string, root string, opts ...func(o *PutOptions)) error

	// GetObject 根据对象路径，获取单个对象
	GetObject(bucket string, key string, opts ...func(o *GetOptions)) (*ObjectOutputHook, error)
	// GetObjectsWithIterator 根据对象前缀，获取多个对象
	GetObjectsWithIterator(bucket string, prefix string, iterator func(keys []*string) error, opts ...func(o *ListOptions)) error

	// DelObject 删除单个对象
	DelObject(bucket string, key string, opts ...func(o *DelOptions)) error
	// DelObjectsWithPrefix 根据对象前缀，删除多个对象
	DelObjectsWithPrefix(bucket string, prefix string, opts ...func(o *DelOptions)) error

	// Copy 对象拷贝
	Copy(src Position, dst Position, opts ...func(o *CopyOptions)) error

	// IsExist 查询对象是否存在
	IsExist(bucket string, key string, opts ...func(o *GetOptions)) (bool, error)

	// Presign 获取对象 HTTP 访问路径
	Presign(bucket string, key string, opts ...func(o *PresignOptions)) (string, error)

	// Downloads 根据前缀，批量下载对象至本地存储，本地存储路径为 'root/prefix/xxx'
	Downloads(bucket string, prefix string, root string, opts ...func(o *ListOptions)) error
}

// C default Client
var C Client

// Init .
func Init(c Client, err error) error {
	C = c
	return err
}

// ObjectInput .
type ObjectInput interface {
	Bucket() string
	Key() string
	ContentType() string
	Close() error
	Error() error
	Body() io.ReadSeeker
}

// input .
type input struct {
	bucket string
	key    string
	ct     content.Type

	err   error
	close func() error
	body  io.ReadSeeker
}

// Bucket .
func (i *input) Bucket() string {
	return i.bucket
}

// Key .
func (i *input) Key() string {
	return i.key
}

// ContentType .
func (i *input) ContentType() string {
	return i.ct.String()
}

// Error .
func (i *input) Error() error {
	return i.err
}

// Close .
func (i *input) Close() error {
	if i.close != nil {
		return i.close()
	} else {
		return nil
	}
}

// Body .
func (i *input) Body() io.ReadSeeker {
	return i.body
}

// InputFile .
func InputFile(bucket string, key string, file string) ObjectInput {
	if f, e := os.Open(file); e != nil {
		return &input{bucket: bucket, key: key, err: e}
	} else {
		return &input{
			bucket: bucket,
			key:    key,
			ct:     content.Stream,
			body:   f,
			close: func() error {
				return f.Close()
			},
		}
	}
}

// InputString .
func InputString(bucket string, key string, v string) ObjectInput {
	return &input{
		bucket: bucket,
		key:    key,
		ct:     content.Text,
		body:   strings.NewReader(v),
	}
}

// InputNumber .
func InputNumber(bucket string, key string, v interface{}) ObjectInput {
	switch v.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return InputString(bucket, key, fmt.Sprintf("%v", v))
	default:
		return &input{bucket: bucket, key: key, err: errors.Errorf(`%T cannot be used as a number.`, v)}
	}
}

// InputBoolean .
func InputBoolean(bucket string, key string, v bool) ObjectInput {
	if v {
		return InputString(bucket, key, "1")
	} else {
		return InputString(bucket, key, "0")
	}
}

// InputMarshalJson .
func InputMarshalJson(bucket string, key string, vPointer interface{}) ObjectInput {
	data, err := json.Marshal(vPointer)
	if err != nil {
		return &input{bucket: bucket, key: key, err: err}
	}
	return &input{
		bucket: bucket,
		key:    key,
		ct:     content.Json,
		body:   bytes.NewReader(data),
	}
}

// InputMarshalProto .
func InputMarshalProto(bucket string, key string, v proto.Message) ObjectInput {
	data, err := proto.Marshal(v)
	if err != nil {
		return &input{bucket: bucket, key: key, err: err}
	}
	return &input{
		bucket: bucket,
		key:    key,
		ct:     content.Proto,
		body:   bytes.NewReader(data),
	}
}

// InputReadSeeker .
func InputReadSeeker(bucket string, key string, body io.ReadSeeker) ObjectInput {
	return &input{
		bucket: bucket,
		key:    key,
		ct:     content.Stream,
		body:   body,
	}
}

// ObjectOutputHook .
type ObjectOutputHook struct {
	// Fetch get object with io.ReadCloser
	Fetch func(hook func(output ObjectOutput) error) error
	// Write write object to io.WriterAt
	Write func(w io.WriterAt) error
}

// ObjectOutput .
type ObjectOutput interface {
	Bucket() string
	Key() string
	ContentType() content.Type
	Decode(vPointer interface{}) error
	Body() io.ReadCloser
	Close() error
}

// output .
type output struct {
	bucket string
	key    string
	ct     string
	body   io.ReadCloser
	once   sync.Once
}

// Decode .
func (o *output) Decode(vPointer interface{}) error {
	buff := new(bytes.Buffer)
	if _, err := io.Copy(buff, o.body); err != nil {
		return err
	}

	o.Close()

	switch vPointer.(type) {
	case *bool:
		if buff.Len() == 1 || buff.Len() == 4 || buff.Len() == 5 {
			switch strings.ToLower(string(buff.Bytes())) {
			case "1", "true":
				*(vPointer.(*bool)) = true
				return nil
			case "0", "false":
				*(vPointer.(*bool)) = false
				return nil
			}
		}
		return errors.New("object decoding failed. incorrect object type.")
	case *[]byte:
		*(vPointer.(*[]byte)) = buff.Bytes()
	case *string:
		*(vPointer.(*string)) = string(buff.Bytes())
	default:
		if pm, ok := vPointer.(proto.Message); ok && o.ct == content.Proto.String() {
			return proto.Unmarshal(buff.Bytes(), pm)
		} else {
			return json.Unmarshal(buff.Bytes(), vPointer)
		}
	}

	return nil
}

// Bucket .
func (o *output) Bucket() string {
	return o.bucket
}

// Key .
func (o *output) Key() string {
	return o.key
}

// ContentType .
func (o *output) ContentType() content.Type {
	return content.Convert(o.ct)
}

// Close .
func (o *output) Close() error {
	o.once.Do(func() {
		o.body.Close()
	})
	return nil
}

// Body .
func (o *output) Body() io.ReadCloser {
	return o.body
}

// OutputReadCloser .
func OutputReadCloser(bucket string, key string, contenttype string, body io.ReadCloser) ObjectOutput {
	return &output{
		bucket: bucket,
		key:    key,
		body:   body,
		ct:     contenttype,
	}
}

// Position .
type Position interface {
	Bucket() string
	Key() string
	IsPrefix() bool
}

// positionRemote .
type positionRemote struct {
	bucket string
	key    string

	isPrefix bool
}

// Bucket .
func (p *positionRemote) Bucket() string {
	return p.bucket
}

// Key .
func (p *positionRemote) Key() string {
	return p.key
}

// IsPrefix .
func (p *positionRemote) IsPrefix() bool {
	return p.isPrefix
}

// PositionRemote .
func PositionRemote(bucket, path string) Position {
	return &positionRemote{
		bucket:   bucket,
		key:      path,
		isPrefix: strings.HasSuffix(path, "/"),
	}
}
