package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"unicode/utf8"

	"google.golang.org/protobuf/proto"

	"github.com/charlesbases/library/content"
	"github.com/charlesbases/library/regexp"
)

var ErrNoSuchKey = errors.New("NoSuchKey: The specified key does not exist.")

var BaseClient Client

type Client interface {
	// PutObject put object to storage
	PutObject(input ObjectInput, opts ...func(o *PutOptions)) error
	// PutFolder put files in folder to storage
	PutFolder(bucket string, prefix string, root string, opts ...func(o *PutOptions)) error

	// GetObject get object
	GetObject(bucket string, key string, opts ...func(o *GetOptions)) (*ObjectOutputHook, error)
	// GetObjectsWithIterator get object list
	GetObjectsWithIterator(bucket string, prefix string, iterator func(keys []*string) error, opts ...func(o *ListOptions)) error

	// DelObject delete object with key
	DelObject(bucket string, key string, opts ...func(o *DelOptions)) error
	// DelPrefix delete object list with prefix
	DelPrefix(bucket string, prefix string, opts ...func(o *DelOptions)) error

	// Copy copy Object to target
	// If the source is prefixed, copy all the objects
	Copy(src Position, dst Position, opts ...func(o *CopyOptions)) error

	// IsExist query whether the object exists
	// If the query is prefixed, the key needs to end with '/'
	IsExist(bucket string, key string, opts ...func(o *GetOptions)) (bool, error)

	// Presign url of object
	Presign(bucket string, key string, opts ...func(o *PresignOptions)) (string, error)

	// Compress compress object into '*.tar.gz'
	// If compressing multiple objects, the key needs to end with '/'
	Compress(bucket string, key string, dst io.Writer, opts ...func(o *ListOptions)) error
}

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

func (i *input) Bucket() string {
	return i.bucket
}

func (i *input) Key() string {
	return i.key
}

func (i *input) ContentType() string {
	return i.ct.String()
}

func (i *input) Error() error {
	return i.err
}

func (i *input) Close() error {
	if i.close != nil {
		return i.close()
	} else {
		return nil
	}
}

func (i *input) Body() io.ReadSeeker {
	return i.body
}

// InputFile .
func InputFile(bucket string, key string, file string) ObjectInput {
	if f, e := os.Open(file); e != nil {
		return &input{err: e}
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
		return &input{err: fmt.Errorf(`%T cannot be used as a number.`, v)}
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
		return &input{err: err}
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
		return &input{err: err}
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

type ObjectOutputHook struct {
	// Fetch get object with io.ReadCloser
	Fetch func(hook func(output ObjectOutput) error) error
	// Write write object to io.WriterAt
	Write func(w io.WriterAt) error
}

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

func (o *output) Bucket() string {
	return o.bucket
}

func (o *output) Key() string {
	return o.key
}

func (o *output) ContentType() content.Type {
	return content.Convert(o.ct)
}

func (o *output) Close() error {
	o.once.Do(func() {
		o.body.Close()
	})
	return nil
}

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

type Position interface {
	Bucket() string
	Path() string
	IsPrefix() bool
}

// positionRemote .
type positionRemote struct {
	bucket string
	path   string

	isPrefix bool
}

func (p *positionRemote) Bucket() string {
	return p.bucket
}

func (p *positionRemote) Path() string {
	return p.path
}

func (p *positionRemote) IsPrefix() bool {
	return p.isPrefix
}

// PositionRemote .
func PositionRemote(bucket, path string) Position {
	return &positionRemote{
		bucket:   bucket,
		path:     path,
		isPrefix: strings.HasSuffix(path, "/"),
	}
}

// CheckBucketName check the compliance of the bucket name.
func CheckBucketName(v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return errors.New("bucket name cannot be empty")
	}
	if regexp.IP.MatchString(v) {
		return errors.New("bucket name cannot be an ip address")
	}

	return nil
}

// CheckObjectNameV1 check the compliance of the object name.
func CheckObjectNameV1(v string) error {
	if len(strings.TrimSpace(v)) == 0 {
		return errors.New("object name cannot be empty")
	}
	if !utf8.ValidString(v) {
		return errors.New("object name with non UTF-8 strings are not supported")
	}
	return nil
}

// CheckObjectNameV2 check the compliance of the object name with CheckObjectNameV1 and verify if the object name ends with '/'.
func CheckObjectNameV2(v string) error {
	if strings.HasSuffix(v, "/") {
		return errors.New("object name cannot end with '/'")
	}
	return CheckObjectNameV1(v)
}
