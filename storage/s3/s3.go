package s3

import (
	"context"
	"crypto/tls"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/charlesbases/salmon"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library/rootpath"
	"github.com/charlesbases/library/storage"
)

const (
	// defaultPoolSize 批量对象操作时的并发数
	defaultPoolSize = 1000
	// defaultFilePerm default of file perm
	defaultFilePerm = 0644
	// defaultFolderPerm default folder perm
	defaultFolderPerm = 0755
	// defaultS3MaxKeys default maxkey with aws-s3
	defaultS3MaxKeys int = 1000
	// defaultS3ManagerPartSize s3manager.Uploader and s3manager.Downloader default PsrtSize. 128 Mib
	defaultS3ManagerPartSize = 128 * 1 << 20
	// defaultS3ManagerConcurrency s3manager.Uploader and s3manager.Downloader default Concurrency.
	defaultS3ManagerConcurrency = 1
)

// s3Client .
type s3Client struct {
	s3   *s3.S3
	opts *storage.Options

	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// pathjoin .
func pathjoin(v ...string) string {
	return strings.Join(v, "/")
}

// awserror .
func awserror(err error) error {
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			return errors.Errorf("%s: %s", awsErr.Code(), awsErr.Message())
		}
	}
	return err
}

// upload use s3manager.Uploader
func (c *s3Client) upload(ctx context.Context, input *s3manager.UploadInput) error {
	_, err := c.uploader.UploadWithContext(ctx, input)
	return awserror(err)
}

// putobject .
func (c *s3Client) putobject(input storage.ObjectInput, opts ...func(o *storage.PutOptions)) error {
	if err := storage.ErrorValidator(
		storage.ValidatorFunc(input.Error),
		storage.BucketName(input.Bucket()),
		storage.KeyName(input.Key()),
	); err != nil {
		return err
	}

	var popts = storage.NewPutOptions(opts...)

	logger.CallerSkip(popts.CallerSkip+1).WithContext(popts.Context).Debugf("[aws-s3]: put(%s.%s)", input.Bucket(), input.Key())

	// close body
	defer input.Close()

	return c.upload(popts.Context, &s3manager.UploadInput{
		Key:         aws.String(input.Key()),
		Bucket:      aws.String(input.Bucket()),
		Body:        input.Body(),
		ContentType: aws.String(input.ContentType()),
	})
}

// PutObject .
func (c *s3Client) PutObject(input storage.ObjectInput, opts ...func(o *storage.PutOptions)) error {
	return errors.Wrapf(c.putobject(input, opts...), "[aws-s3]: put(%s.%s)", input.Bucket(), input.Key())
}

// putfolder .
func (c *s3Client) putfolder(bucket string, prefix string, root string, opts ...func(o *storage.PutOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyPrefixName(prefix)); err != nil {
		return err
	}

	r := rootpath.NewRoot(root)
	if !r.IsDir() {
		return errors.Errorf(`"%s" is not a folder.`, root)
	}

	var popts = storage.NewPutOptions(opts...)

	logger.CallerSkip(popts.CallerSkip+1).WithContext(popts.Context).Debugf("[aws-s3]: put(%s.%s.*)", bucket, prefix)

	pool, err := salmon.NewPool(defaultPoolSize, func(v interface{}, stop func()) {
		if key, ok := v.(*string); ok {
			input := storage.InputFile(bucket, filepath.ToSlash(filepath.Join(prefix, strings.TrimPrefix(*key, root))), *key)
			if err := input.Error(); err != nil {
				stop()
				logger.WithContext(popts.Context).Errorf("[aws-s3]: put(%s.%s.*): %v", bucket, *key, err)
			} else {
				if err := c.upload(popts.Context, &s3manager.UploadInput{
					Key:         aws.String(input.Key()),
					Bucket:      aws.String(input.Bucket()),
					Body:        input.Body(),
					ContentType: aws.String(input.ContentType()),
				}); err != nil {
					stop()
					logger.WithContext(popts.Context).Errorf("[aws-s3]: put(%s.%s.*): %v", bucket, *key, err)
				}

				// close body
				input.Close()
			}
		}
	})
	if err != nil {
		return err
	}
	defer pool.Wait()

	return r.Walk(func(path string, info fs.FileInfo) error {
		if !info.IsDir() {
			pool.Invoke(&path)
		}
		return nil
	})
}

// PutFolder .
func (c *s3Client) PutFolder(bucket string, prefix string, root string, opts ...func(o *storage.PutOptions)) error {
	return errors.Wrapf(c.putfolder(bucket, prefix, root, opts...), "[aws-s3]: put(%s.%s.*)", bucket, prefix)
}

// download use s3manager.Downloader
func (c *s3Client) download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput) error {
	_, err := c.downloader.DownloadWithContext(ctx, w, input)
	return awserror(err)
}

// GetObject .
func (c *s3Client) GetObject(bucket string, key string, opts ...func(o *storage.GetOptions)) (*storage.ObjectOutputHook, error) {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyName(key)); err != nil {
		return nil, errors.Wrapf(err, "[aws-s3]: get(%s.%s)", bucket, key)
	}

	var gopts = storage.NewGetOptions(opts...)

	logger.CallerSkip(gopts.CallerSkip).WithContext(gopts.Context).Debugf("[aws-s3]: get(%s.%s)", bucket, key)

	input := &s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(gopts.VersionID),
	}

	return &storage.ObjectOutputHook{
		Fetch: func(hook func(output storage.ObjectOutput) error) error {
			return errors.Wrapf(func() error {
				output, err := c.s3.GetObjectWithContext(gopts.Context, input)
				if err != nil {
					return awserror(err)
				}
				return hook(storage.OutputReadCloser(bucket, key, aws.StringValue(output.ContentType), output.Body))
			}(), "[aws-s3]: get(%s.%s)", bucket, key)
		},
		Write: func(w io.WriterAt) error {
			return errors.Wrapf(c.download(gopts.Context, w, input), "[aws-s3]: get(%s.%s)", bucket, key)
		},
	}, nil
}

// listobjects .
func (c *s3Client) listobjects(bucket string, prefix string, lopts *storage.ListOptions, iterator func(keys []*string) error) error {
	// -1 < lopts.MaxKeys < defaultS3MaxKeys
	var offset = defaultS3MaxKeys
	if lopts.MaxKeys > -1 && lopts.MaxKeys < offset {
		offset = lopts.MaxKeys
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		MaxKeys:   aws.Int64(int64(offset)),
		Delimiter: aws.String("/"),
	}

	if lopts.Recursive {
		// If recursive we do not delimit.
		input.Delimiter = nil
	}

	var count int

	return awserror(c.s3.ListObjectsV2PagesWithContext(lopts.Context, input,
		func(output *s3.ListObjectsV2Output, lasted bool) bool {
			if len(output.Contents) != 0 {
				keys := make([]*string, 0, len(output.Contents))
				for _, content := range output.Contents {
					keys = append(keys, content.Key)
				}

				// do something
				if err := iterator(keys); err != nil {
					return false
				}

				if lopts.MaxKeys > 0 {
					count += len(keys)

					if (lopts.MaxKeys - count) < offset {
						input.MaxKeys = aws.Int64(int64(lopts.MaxKeys - count))
					}
				}
			}

			// listing ends result is not truncated, return right here.
			if lasted {
				return false
			}

			return count != lopts.MaxKeys
		}))
}

// GetObjectsWithIterator .
func (c *s3Client) GetObjectsWithIterator(bucket string, prefix string, iterator func(keys []*string) error, opts ...func(o *storage.ListOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyPrefixName(prefix)); err != nil {
		return errors.Wrapf(err, "[aws-s3]: get(%s.%s.*)", bucket, prefix)
	}

	var lopts = storage.NewListOptions(opts...)

	logger.CallerSkip(lopts.CallerSkip).WithContext(lopts.Context).Debugf("[aws-s3]: get(%s.%s.*)", bucket, prefix)

	return errors.Wrapf(c.listobjects(bucket, prefix, lopts, iterator), "[aws-s3]: get(%s.%s.*)", bucket, prefix)
}

// DelObject .
func (c *s3Client) DelObject(bucket string, key string, opts ...func(o *storage.DelOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyName(key)); err != nil {
		return errors.Wrapf(err, "[aws-s3]: del(%s.%s)", bucket, key)
	}

	var dopts = storage.NewDelOptions(opts...)

	logger.CallerSkip(dopts.CallerSkip).WithContext(dopts.Context).Debugf("[aws-s3]: del(%s.%s)", bucket, key)

	_, err := c.s3.DeleteObjectWithContext(dopts.Context, &s3.DeleteObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(dopts.VersionID),
	})
	return errors.Wrapf(awserror(err), "[aws-s3]: del(%s.%s)", bucket, key)
}

// DelObjectsWithPrefix .
func (c *s3Client) DelObjectsWithPrefix(bucket string, prefix string, opts ...func(o *storage.DelOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyPrefixName(prefix)); err != nil {
		return errors.Wrapf(err, "[aws-s3]: del(%s.%s.*)", bucket, prefix)
	}

	var dopts = storage.NewDelOptions(opts...)

	logger.CallerSkip(dopts.CallerSkip).WithContext(dopts.Context).Debugf("[aws-s3]: del(%s.%s.*)", bucket, prefix)

	pool, err := salmon.NewPool(defaultPoolSize, func(v interface{}, stop func()) {
		if keys, ok := v.([]*string); ok {
			var items = make([]*s3.ObjectIdentifier, 0, len(keys))
			for _, key := range keys {
				items = append(items, &s3.ObjectIdentifier{
					Key: key,
				})
			}

			_, err := c.s3.DeleteObjectsWithContext(dopts.Context, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3.Delete{
					Objects: items,
					Quiet:   aws.Bool(true),
				},
			})
			if err != nil {
				stop()
				logger.WithContext(dopts.Context).Errorf("[aws-s3]: del(%s.%s.*): %v", bucket, prefix, awserror(err))
			}
		}
	})
	if err != nil {
		return errors.Wrapf(err, "[aws-s3]: del(%s.%s.*)", bucket, prefix)
	}
	defer pool.Wait()

	return errors.Wrapf(c.listobjects(bucket, prefix,
		storage.NewListOptions(func(o *storage.ListOptions) {
			o.Context = dopts.Context
			o.MaxKeys = -1
			o.Recursive = true
		}),
		func(keys []*string) error {
			var length = len(keys)

			// 一个协程只处理十个 key
			left, right := 0, 0
			for left != length {
				if right = left + 10; right >= length {
					right = length
				}

				pool.Invoke(keys[left:right])
				left = right
			}
			return nil
		},
	), "[aws-s3]: del(%s.%s.*)", bucket, prefix)
}

// copyobject .
func (c *s3Client) copyobject(ctx context.Context, input *s3.CopyObjectInput) error {
	_, err := c.s3.CopyObjectWithContext(ctx, input)
	return awserror(err)
}

// copy .
func (c *s3Client) copy(src, dst storage.Position, opts ...func(o *storage.CopyOptions)) error {
	if err := storage.ErrorValidator(
		storage.BucketName(src.Bucket()),
		storage.BucketName(dst.Bucket()),
		storage.KeyPrefixName(src.Key()),
		storage.KeyPrefixName(dst.Key()),
		storage.ValidatorFunc(func() error {
			if src.IsPrefix() != dst.IsPrefix() {
				return errors.New("the src and dst position must either end with '/', or both should be keys")
			}
			return nil
		})); err != nil {
		return err
	}

	var copts = storage.NewCopyOptions(opts...)

	switch src.IsPrefix() {
	case false:
		logger.CallerSkip(copts.CallerSkip+1).WithContext(copts.Context).Debugf(`[aws-s3]: copy("%s.%s" -> "%s.%s")`, src.Bucket(), src.Key(), dst.Bucket(), dst.Key())

		return c.copyobject(copts.Context,
			&s3.CopyObjectInput{
				CopySource: aws.String(pathjoin(src.Bucket(), src.Key())),
				Bucket:     aws.String(dst.Bucket()),
				Key:        aws.String(dst.Key()),
			})
	default:
		logger.CallerSkip(copts.CallerSkip+1).WithContext(copts.Context).Debugf(`[aws-s3] copy("%s.%s.*" -> "%s.%s.*")`, src.Bucket(), src.Key(), dst.Bucket(), dst.Key())

		pool, err := salmon.NewPool(defaultPoolSize, func(v interface{}, stop func()) {
			if keys, ok := v.([]*string); ok {
				for _, key := range keys {
					if err := c.copyobject(copts.Context,
						&s3.CopyObjectInput{
							CopySource: aws.String(pathjoin(src.Bucket(), *key)),
							Bucket:     aws.String(dst.Bucket()),
							Key:        aws.String(pathjoin(dst.Key(), strings.TrimPrefix(*key, src.Key()))),
						}); err != nil {
					}
				}
			}
		})
		if err != nil {
			return err
		}
		defer pool.Wait()

		return c.listobjects(src.Bucket(), src.Key(),
			storage.NewListOptions(func(o *storage.ListOptions) {
				o.Context = copts.Context
				o.MaxKeys = -1
				o.Recursive = true
			}),
			func(keys []*string) error {
				var length = len(keys)

				// 一个协程只处理十个 key
				left, right := 0, 0
				for left != length {
					if right = left + 10; right >= length {
						right = length
					}

					pool.Invoke(keys[left:right])
					left = right
				}
				return nil
			})
	}
}

// Copy .
func (c *s3Client) Copy(src, dst storage.Position, opts ...func(o *storage.CopyOptions)) error {
	return errors.Wrapf(c.copy(src, dst, opts...), `[aws-s3]: copy("%s.%s" -> "%s.%s")`, src.Bucket(), src.Key(), dst.Bucket(), dst.Key())
}

// headObject .
func (c *s3Client) headObject(bucket string, key string, gopts *storage.GetOptions) (*s3.HeadObjectOutput, error) {
	head, err := c.s3.HeadObjectWithContext(gopts.Context, &s3.HeadObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(gopts.VersionID),
	})
	return head, awserror(err)
}

// exists .
func (c *s3Client) exists(bucket, key string, opts ...func(o *storage.GetOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyPrefixName(key)); err != nil {
		return err
	}

	var gopts = storage.NewGetOptions(opts...)

	switch strings.HasSuffix(key, "/") {
	// object
	case false:
		_, err := c.headObject(bucket, key, gopts)
		return err
	// prefix
	default:
		var isExist bool
		err := c.listobjects(bucket, key,
			storage.NewListOptions(func(o *storage.ListOptions) {
				o.Context = gopts.Context
				o.MaxKeys = 1
				o.Recursive = true
			}),
			func(keys []*string) error {
				isExist = true
				return nil
			},
		)
		if !isExist {
			return storage.ErrNoSuchKey
		}
		return err
	}
}

// IsExist .
func (c *s3Client) IsExist(bucket, key string, opts ...func(o *storage.GetOptions)) (bool, error) {
	err := c.exists(bucket, key, opts...)
	return err == nil, errors.Wrapf(err, "[aws-s3]: exists(%s.%s)", bucket, key)
}

// presign .
func (c *s3Client) presign(bucket, key string, opts ...func(o *storage.PresignOptions)) (string, error) {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyName(key)); err != nil {
		return "", err
	}

	var popts = storage.NewPresignOptions(opts...)

	request, _ := c.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(popts.VersionID),
	})

	url, err := request.Presign(popts.Expires)
	return url, awserror(err)
}

// Presign .
func (c *s3Client) Presign(bucket, key string, opts ...func(o *storage.PresignOptions)) (string, error) {
	url, err := c.presign(bucket, key, opts...)
	return url, errors.Wrapf(err, "[aws-s3]: presign(%s.%s)", bucket, key)
}

// downloads 。
func (c *s3Client) downloads(bucket string, prefix string, root string, opts ...func(o *storage.ListOptions)) error {
	if err := storage.ErrorValidator(storage.BucketName(bucket), storage.KeyPrefixName(prefix)); err != nil {
		return err
	}

	var lopts = storage.NewListOptions(opts...)

	pool, err := salmon.NewPool(defaultPoolSize, func(v interface{}, stop func()) {
		if keys, ok := v.([]*string); ok {
			for _, objkey := range keys {
				if err := func() error {
					file, err := openfile(filepath.Join(root, strings.Replace(aws.StringValue(objkey), prefix, "./", 1)))
					if err != nil {
						return err
					}
					defer file.Close()

					return c.download(lopts.Context, file,
						&s3.GetObjectInput{
							Bucket: aws.String(bucket),
							Key:    objkey,
						})
				}(); err != nil {
					stop()
					logger.WithContext(lopts.Context).Errorf("[aws-s3]: downloads(%s.%s.*): %v", bucket, *objkey, err)
				}
			}
		}
	})
	if err != nil {
		return err
	}
	defer pool.Wait()

	logger.CallerSkip(lopts.CallerSkip+1).WithContext(lopts.Context).Debugf("[aws-s3]: downloads(%s.%s.*)", bucket, prefix)

	return c.listobjects(bucket, prefix,
		storage.NewListOptions(func(o *storage.ListOptions) {
			o.Context = lopts.Context
			o.MaxKeys = -1
			o.Recursive = true
		}),
		func(keys []*string) error {
			var length = len(keys)

			// 一个协程只处理十个 key
			left, right := 0, 0
			for left != length {
				if right = left + 10; right >= length {
					right = length
				}

				pool.Invoke(keys[left:right])
				left = right
			}
			return nil
		},
	)
}

// Downloads .
func (c *s3Client) Downloads(bucket string, prefix string, root string, opts ...func(o *storage.ListOptions)) error {
	return errors.Wrapf(c.downloads(bucket, prefix, root, opts...), "[aws-s3]: downloads(%s.%s.*)", bucket, prefix)
}

// ping .
func (c *s3Client) ping() error {
	_, err := c.s3.ListBuckets(nil)
	return errors.Wrapf(awserror(err), "[aws-s3]: ping")
}

// openfile .
func openfile(name string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(name), defaultFolderPerm); err != nil {
		return nil, err
	}
	return os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, defaultFilePerm)
}

// NewClient .
func NewClient(endpoint string, accessKey string, secretKey string, opts ...func(o *storage.Options)) (storage.Client, error) {
	client := &s3Client{opts: storage.NewOptions(opts...)}

	// new client
	sess, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(client.opts.Region),
		DisableSSL:       aws.Bool(!client.opts.UseSSL),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		HTTPClient: &http.Client{
			Timeout: client.opts.Timeout,
			Transport: &http.Transport{
				// DisableKeepAlives: false,
				// MaxIdleConns:      1 << 10,
				// IdleConnTimeout:   time.Second * 30,
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	client.s3 = s3.New(sess)

	client.uploader = s3manager.NewUploaderWithClient(client.s3, func(uploader *s3manager.Uploader) {
		uploader.PartSize = defaultS3ManagerPartSize
		uploader.Concurrency = defaultS3ManagerConcurrency
	})
	client.downloader = s3manager.NewDownloaderWithClient(client.s3, func(downloader *s3manager.Downloader) {
		downloader.PartSize = defaultS3ManagerPartSize
		downloader.Concurrency = defaultS3ManagerConcurrency
	})

	return client, client.ping()
}
