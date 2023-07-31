package s3

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/charlesbases/logger"

	"github.com/charlesbases/library/compress/gzip"
	"github.com/charlesbases/library/rootpath"
	"github.com/charlesbases/library/storage"
)

const (
	// defaultS3MaxKeys default maxkey with aws-s3
	defaultS3MaxKeys int = 1000
	// defaultPartSize s3manager.Uploader and s3manager.Downloader default PsrtSize. 128 Mib
	defaultPartSize = 128 * 1 << 20
	// defaultPartSize s3manager.Uploader and s3manager.Downloader default Concurrency.
	// 压缩文件时，archive/tar 不支持 io.WaritAt 随机写入，所以不可使用协程并发写入，只能顺序写入。
	defaultConcurrency = 1
)

// s3Client .
type s3Client struct {
	s3 *s3.S3

	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader

	options *storage.Options
}

// pathJoin .
func pathJoin(v ...string) string {
	return strings.Join(v, "/")
}

// decodeAwsErr .
func decodeAwsErr(err error) error {
	if awsErr, ok := err.(awserr.Error); ok {
		return fmt.Errorf("%s: %s", awsErr.Code(), awsErr.Message())
	}
	return err
}

// upload use s3manager.Uploader
func (client *s3Client) upload(ctx context.Context, input *s3manager.UploadInput) error {
	_, err := client.uploader.UploadWithContext(ctx, input)
	if err != nil {
		err = decodeAwsErr(err)
		logger.Errorf("[aws-s3] put(%s.%s) failed. '%s'", aws.StringValue(input.Bucket), aws.StringValue(input.Key), err.Error())
		return err
	}
	return nil
}

func (client *s3Client) PutObject(input storage.ObjectInput, opts ...func(o *storage.PutOptions)) error {
	var popts = storage.NewPutOptions(opts...)

	if err := input.Error(); err != nil {
		return err
	}
	if err := storage.CheckBucketName(input.Bucket()); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV2(input.Key()); err != nil {
		return err
	}

	logger.Debugf("[aws-s3] put(%s.%s)", input.Bucket(), input.Key())

	// close body
	defer input.Close()

	return client.upload(popts.Context, &s3manager.UploadInput{
		Key:         aws.String(input.Key()),
		Bucket:      aws.String(input.Bucket()),
		Body:        input.Body(),
		ContentType: aws.String(input.ContentType()),
	})
}

func (client *s3Client) PutFolder(bucket string, prefix string, root string, opts ...func(o *storage.PutOptions)) error {
	var popts = storage.NewPutOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(prefix); err != nil {
		return err
	}

	r := rootpath.NewRoot(root)
	if !r.IsDir() {
		return fmt.Errorf(`"%s" is not a folder.`, root)
	}

	if !strings.HasSuffix(prefix, "/") {
		return errors.New("the prefix must end in '/'.")
	}

	logger.Debugf("[aws-s3] put(%s.%s.*)", bucket, prefix)

	f := storage.NewFactory()

	// putobject
	f.Flowline(func(v interface{}) {
		if key, ok := v.(*string); ok {
			input := storage.InputFile(bucket, filepath.ToSlash(filepath.Join(prefix, strings.TrimPrefix(*key, root))), *key)
			if err := input.Error(); err != nil {
				logger.Errorf("[aws-s3] put(%s.%s.*)", input.Bucket(), input.Key(), err.Error())
			} else {
				client.upload(popts.Context, &s3manager.UploadInput{
					Key:         aws.String(input.Key()),
					Bucket:      aws.String(input.Bucket()),
					Body:        input.Body(),
					ContentType: aws.String(input.ContentType()),
				})

				// close body
				input.Close()
			}
		}
	})

	err := r.Walk(func(path string, info fs.FileInfo) error {
		if !info.IsDir() {
			f.Push(&path)
		}
		return nil
	})

	f.Wait()
	return err
}

// download use s3manager.Downloader
func (client *s3Client) download(ctx context.Context, w io.WriterAt, input *s3.GetObjectInput) error {
	_, err := client.downloader.DownloadWithContext(ctx, w, input)
	if err != nil {
		err = decodeAwsErr(err)
		logger.Errorf("[aws-s3] get(%s.%s) failed. '%s'", aws.StringValue(input.Bucket), aws.StringValue(input.Key), err.Error())
		return err
	}
	return nil
}

func (client *s3Client) GetObject(bucket string, key string, opts ...func(o *storage.GetOptions)) (*storage.ObjectOutputHook, error) {
	var gopts = storage.NewGetOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return nil, err
	}
	if err := storage.CheckObjectNameV2(key); err != nil {
		return nil, err
	}

	logger.Debugf("[aws-s3] get(%s.%s)", bucket, key)

	input := &s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(gopts.VersionID),
	}

	return &storage.ObjectOutputHook{
		Fetch: func(hook func(output storage.ObjectOutput) error) error {
			s3Output, err := client.s3.GetObjectWithContext(gopts.Context, input)
			if err != nil {
				err = decodeAwsErr(err)
				logger.Errorf("[aws-s3] get(%s.%s) failed. '%s'", bucket, key, err.Error())
				return err
			}
			return hook(storage.OutputReadCloser(bucket, key, aws.StringValue(s3Output.ContentType), s3Output.Body))
		},
		Write: func(w io.WriterAt) error {
			return client.download(gopts.Context, w, input)
		},
	}, nil
}

// listobjects .
func (client *s3Client) listobjects(bucket string, prefix string, lopts *storage.ListOptions, iterator func(keys []*string) error) error {
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

	return client.s3.ListObjectsV2PagesWithContext(lopts.Context, input,
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
		})
}

func (client *s3Client) GetObjectsWithIterator(bucket string, prefix string, iterator func(keys []*string) error, opts ...func(o *storage.ListOptions)) error {
	if err := storage.CheckBucketName(bucket); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(prefix); err != nil {
		return err
	}

	logger.Debugf("[aws-s3] get(%s.%s.*)", bucket, prefix)

	if err := client.listobjects(bucket, prefix, storage.NewListOptions(opts...), iterator); err != nil {
		err = decodeAwsErr(err)
		logger.Errorf("[aws-s3] get(%s.%s.*) failed. '%s'", bucket, prefix, err.Error())
		return err
	}
	return nil
}

func (client *s3Client) DelObject(bucket string, key string, opts ...func(o *storage.DelOptions)) error {
	var dopts = storage.NewDelOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV2(key); err != nil {
		return err
	}

	logger.Debugf("[aws-s3] del(%s.%s)", bucket, key)

	if _, err := client.s3.DeleteObjectWithContext(dopts.Context, &s3.DeleteObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(dopts.VersionID),
	}); err != nil {
		err = decodeAwsErr(err)
		logger.Errorf("[aws-s3] del(%s.%s)[key] failed. '%s'", bucket, key, err.Error())
		return err
	}
	return nil
}

func (client *s3Client) DelPrefix(bucket string, prefix string, opts ...func(o *storage.DelOptions)) error {
	var dopts = storage.NewDelOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(prefix); err != nil {
		return err
	}

	logger.Debugf("[aws-s3] del(%s.%s.*)", bucket, prefix)

	f := storage.NewFactory()
	f.Flowline(func(v interface{}) {
		if keys, ok := v.([]*string); ok {
			var items = make([]*s3.ObjectIdentifier, 0, len(keys))
			for _, key := range keys {
				items = append(items, &s3.ObjectIdentifier{
					Key: key,
				})
			}

			_, err := client.s3.DeleteObjectsWithContext(dopts.Context, &s3.DeleteObjectsInput{
				Bucket: aws.String(bucket),
				Delete: &s3.Delete{
					Objects: items,
					Quiet:   aws.Bool(true),
				},
			})
			if err != nil {
				f.Closing()
				logger.Errorf("[aws-s3] del(%s.%s.*) failed. '%s'", bucket, prefix, decodeAwsErr(err).Error())
			}
		}
	})

	err := client.listobjects(bucket, prefix,
		storage.NewListOptions(func(o *storage.ListOptions) {
			o.Context = dopts.Context
			o.MaxKeys = -1
			o.Recursive = true
		}),
		func(keys []*string) error {
			f.PushSlice(keys)
			return nil
		},
	)
	if err != nil {
		f.Closing()
		logger.Errorf("[aws-s3] del(%s.%s.*) failed. '%s'", bucket, prefix, decodeAwsErr(err).Error())
	}

	f.Wait()
	return nil
}

// copyobject .
func (client *s3Client) copyobject(ctx context.Context, srcBucket, srcKey string, input *s3.CopyObjectInput) error {
	if _, err := client.s3.CopyObjectWithContext(ctx, input); err != nil {
		err = decodeAwsErr(err)
		logger.Errorf(`[aws-s3] copy("%s.%s" -> "%s.%s") failed. "%s"`, srcBucket, srcKey, aws.StringValue(input.Bucket), aws.StringValue(input.Key), err.Error())
		return err
	}
	return nil
}

func (client *s3Client) Copy(src, dst storage.Position, opts ...func(o *storage.CopyOptions)) error {
	var copts = storage.NewCopyOptions(opts...)

	if err := storage.CheckBucketName(src.Bucket()); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(src.Path()); err != nil {
		return err
	}
	if err := storage.CheckBucketName(dst.Bucket()); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(dst.Path()); err != nil {
		return err
	}

	if src.IsPrefix() != dst.IsPrefix() {
		return errors.New("the source and destination position must either end with '/', or both should be keys")
	}

	switch src.IsPrefix() {
	case false:
		logger.Debugf(`[aws-s3] copy("%s.%s" -> "%s.%s")`, src.Bucket(), src.Path(), dst.Bucket(), dst.Path())

		return client.copyobject(copts.Context, src.Bucket(), src.Path(),
			&s3.CopyObjectInput{
				CopySource: aws.String(pathJoin(src.Bucket(), src.Path())),
				Bucket:     aws.String(dst.Bucket()),
				Key:        aws.String(dst.Path()),
			})
	default:
		logger.Debugf(`[aws-s3] copy("%s.%s.*" -> "%s.%s.*")`, src.Bucket(), src.Path(), dst.Bucket(), dst.Path())

		f := storage.NewFactory()
		f.Flowline(func(v interface{}) {
			if keys, ok := v.([]*string); ok {
				for _, key := range keys {
					if err := client.copyobject(copts.Context, src.Bucket(), *key,
						&s3.CopyObjectInput{
							CopySource: aws.String(pathJoin(src.Bucket(), *key)),
							Bucket:     aws.String(dst.Bucket()),
							Key:        aws.String(pathJoin(dst.Path(), strings.TrimPrefix(*key, src.Path()))),
						}); err != nil {
					}
				}
			}
		})
		defer f.Wait()

		if err := client.listobjects(src.Bucket(), src.Path(),
			storage.NewListOptions(func(o *storage.ListOptions) {
				o.Context = copts.Context
				o.MaxKeys = -1
				o.Recursive = true
			}),
			func(keys []*string) error {
				f.PushSlice(keys)
				return nil
			}); err != nil {
			f.Closing()
			err = decodeAwsErr(err)
			logger.Errorf(`[aws-s3] copy("%s.%s.*" -> "%s.%s.*") failed. "%s"`, src.Bucket(), src.Path(), dst.Bucket(), dst.Path(), err.Error())
			return err
		}
		return nil
	}
}

// headObject .
func (client *s3Client) headObject(bucket string, key string, gopts *storage.GetOptions) (*s3.HeadObjectOutput, error) {
	head, err := client.s3.HeadObjectWithContext(gopts.Context, &s3.HeadObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(gopts.VersionID),
	})
	if err != nil {
		// not found
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == 404 {
			return nil, nil
		}
		// others error
		return nil, err
	}
	return head, nil
}

func (client *s3Client) IsExist(bucket, key string, opts ...func(o *storage.GetOptions)) (bool, error) {
	var gopts = storage.NewGetOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return false, err
	}
	if err := storage.CheckObjectNameV1(key); err != nil {
		return false, err
	}

	switch strings.HasSuffix(key, "/") {
	// object
	case false:
		head, err := client.headObject(bucket, key, gopts)
		// head !=nil && err == nil
		if head != nil {
			return true, nil
		}
		return false, err
	// prefix
	default:
		var isExist bool

		err := client.listobjects(bucket, key,
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

		return isExist, decodeAwsErr(err)
	}
}

func (client *s3Client) Presign(bucket, key string, opts ...func(o *storage.PresignOptions)) (string, error) {
	var popts = storage.NewPresignOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return "", err
	}
	if err := storage.CheckObjectNameV2(key); err != nil {
		return "", err
	}

	request, _ := client.s3.GetObjectRequest(&s3.GetObjectInput{
		Bucket:    aws.String(bucket),
		Key:       aws.String(key),
		VersionId: aws.String(popts.VersionID),
	})

	if url, err := request.Presign(popts.Expires); err != nil {
		return "", decodeAwsErr(err)
	} else {
		return url, nil
	}
}

func (client *s3Client) Compress(bucket string, key string, dst io.Writer, opts ...func(o *storage.ListOptions)) error {
	var lopts = storage.NewListOptions(opts...)

	if err := storage.CheckBucketName(bucket); err != nil {
		return err
	}
	if err := storage.CheckObjectNameV1(key); err != nil {
		return err
	}

	tx := gzip.New(dst)
	defer tx.Close()

	switch strings.HasSuffix(key, "/") {
	// object
	case false:
		logger.Debugf("[aws-s3] compress(%s.%s)", bucket, key)

		head, err := client.headObject(bucket, key, &storage.GetOptions{Context: lopts.Context})
		if err != nil {
			err = decodeAwsErr(err)
			logger.Errorf("[aws-s3] compress(%s.%s) failed. '%s'", bucket, key, err.Error())
			return err
		}
		if head == nil {
			err = storage.ErrNoSuchKey
			logger.Errorf("[aws-s3] compress(%s.%s) failed. '%s'", bucket, key, err.Error())
			return err
		}

		return tx.WithWriteAt(func(at io.WriterAt) error {
			if err := client.download(lopts.Context, at,
				&s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key:    aws.String(key),
				}); err != nil {
				err = decodeAwsErr(err)
				logger.Errorf("[aws-s3] compress(%s.%s) failed. '%s'", bucket, key, err.Error())
				return err
			}
			return nil
		},
			func(h *gzip.Header) {
				h.Name = filepath.Base(key)
				h.Size = aws.Int64Value(head.ContentLength)
			})
	// prefix
	default:
		logger.Debugf("[aws-s3] compress(%s.%s.*)", bucket, key)

		f := storage.NewFactory()
		f.Flowline(func(v interface{}) {
			if keys, ok := v.([]*string); ok {
				for _, objkey := range keys {
					head, err := client.headObject(bucket, aws.StringValue(objkey), &storage.GetOptions{Context: lopts.Context})
					if err != nil {
						err = decodeAwsErr(err)
						logger.Errorf("[aws-s3] compress(%s.%s) failed. '%s'", bucket, aws.StringValue(objkey), err.Error())
						return
					}

					tx.WithWriteAt(func(at io.WriterAt) error {
						if err := client.download(lopts.Context, at,
							&s3.GetObjectInput{
								Bucket: aws.String(bucket),
								Key:    objkey,
							}); err != nil {
							f.Closing()
							err = decodeAwsErr(err)
							logger.Errorf("[aws-s3] compress(%s.%s.*) failed. '%s'", bucket, aws.StringValue(objkey), err.Error())
							return err
						}
						return nil
					},
						func(h *gzip.Header) {
							h.Name = strings.Replace(aws.StringValue(objkey), key, "./", 1)
							h.Size = aws.Int64Value(head.ContentLength)
						})
				}
			}
		})

		err := client.listobjects(bucket, key,
			storage.NewListOptions(func(o *storage.ListOptions) {
				o.Context = lopts.Context
				o.MaxKeys = -1
				o.Recursive = true
			}),
			func(keys []*string) error {
				f.PushSlice(keys)
				return nil
			},
		)
		if err != nil {
			f.Closing()
			err = decodeAwsErr(err)
			logger.Debugf("[aws-s3] compress(%s.%s.*) failed. '%s'", bucket, key, err.Error())
		}

		f.Wait()
		return err
	}
}

// NewClient .
func NewClient(endpoint string, accessKey string, secretKey string, opts ...func(o *storage.Options)) (storage.Client, error) {
	client := &s3Client{options: storage.NewOptions(opts...)}

	// new client
	session, err := session.NewSession(&aws.Config{
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(client.options.Region),
		DisableSSL:       aws.Bool(!client.options.UseSSL),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		HTTPClient: &http.Client{
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

	client.s3 = s3.New(session)

	client.uploader = s3manager.NewUploaderWithClient(client.s3, func(uploader *s3manager.Uploader) {
		uploader.PartSize = defaultPartSize
		uploader.Concurrency = defaultConcurrency
	})
	client.downloader = s3manager.NewDownloaderWithClient(client.s3, func(downloader *s3manager.Downloader) {
		downloader.PartSize = defaultPartSize
		downloader.Concurrency = defaultConcurrency
	})

	return client, client.ping()
}

// ping .
func (client *s3Client) ping() error {
	ctx, _ := context.WithTimeout(context.Background(), client.options.Timeout)
	if _, err := client.s3.ListBucketsWithContext(ctx, &s3.ListBucketsInput{}); err != nil {
		err = decodeAwsErr(err)
		logger.Errorf(`[aws-s3] dial "%s" failed. '%s'`, client.s3.Endpoint, err.Error())
		return err
	}
	return nil
}

// writeat .
type writeat struct {
	io.WriteCloser
}

// WriteAt 顺序写入
func (w *writeat) WriteAt(p []byte, off int64) (int, error) {
	return w.Write(p)
}
