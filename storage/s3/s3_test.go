package s3

import (
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"github.com/charlesbases/logger"

	"github.com/charlesbases/library/storage"
)

var (
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

	key = func(v string) string {
		return "testdata/data/" + v
	}

	mxdata = "mxdata"
)

func TestS3Client_PutObject(t *testing.T) {
	cli, _ := clientBOS()

	// string
	{
		if err := cli.PutObject(storage.InputString(mxdata, key("simple/string"), time.Now().String())); err != nil {
			logger.Fatalf("string ==> %s", err.Error())
		}
	}

	// number
	{
		if err := cli.PutObject(storage.InputNumber(mxdata, key("simple/number.int"), time.Now().Second())); err != nil {
			logger.Fatalf("number ==> %s", err.Error())
		}
		if err := cli.PutObject(storage.InputNumber(mxdata, key("simple/number.folat"), math.Pi)); err != nil {
			logger.Fatalf("number ==> %s", err.Error())
		}
		// if err := cli.PutObject(storage.InputNumber(mxdata, key("simple/number.string"), time.Now())); err != nil {
		// 	logger.Errorf("number ==> %s", err.Error())
		// }
	}

	// boolean
	{
		if err := cli.PutObject(storage.InputBoolean(mxdata, key("simple/boolean"), true)); err != nil {
			logger.Fatalf("boolean ==> %s", err.Error())
		}
	}

	// json
	{
		var v = map[string]string{"date": time.Now().String()}
		if err := cli.PutObject(storage.InputMarshalJson(mxdata, key("simple/object.json"), &v)); err != nil {
			logger.Fatalf("object.json ==> %s", err.Error())
		}
	}

	// file
	{
		if err := cli.PutObject(storage.InputFile(mxdata, key("simple/file"), "os.go")); err != nil {
			logger.Fatalf("file ==> %s", err.Error())
		}
	}

	// io.ReadSeeker
	{
		f, err := os.Open("os.go")
		if err != nil {
			logger.Fatal(err.Error())
		}
		defer f.Close()

		if err := cli.PutObject(storage.InputReadSeeker(mxdata, key("simple/readseeker"), f)); err != nil {
			logger.Fatalf("readseeker ==> %s", err.Error())
		}
	}
}

func TestS3Client_GetObject(t *testing.T) {
	cli, _ := clientBOS()

	// string
	{
		if hook, err := cli.GetObject(mxdata, key("strinag")); err != nil {
			logger.Fatalf("string ==> %s", err.Error())
		} else {
			var v string

			if err := hook.Fetch(func(output storage.ObjectOutput) error {
				return output.Decode(&v)
			}); err != nil {
				logger.Fatalf("string ==> %s", err.Error())
			}

			logger.Debugf(`string ==> %s`, v)
		}
	}

	// number
	{
		if hook, err := cli.GetObject(mxdata, key("number.int")); err != nil {
			logger.Fatalf("number.int ==> %s", err.Error())
		} else {
			var v int

			if err := hook.Fetch(func(output storage.ObjectOutput) error {
				return output.Decode(&v)
			}); err != nil {
				logger.Fatalf("number.int ==> %s", err.Error())
			}

			logger.Debugf(`number.int ==> %d`, v)
		}

		if hook, err := cli.GetObject(mxdata, key("number.folat")); err != nil {
			logger.Fatalf("number.folat ==> %s", err.Error())
		} else {
			var v float64

			if err := hook.Fetch(func(output storage.ObjectOutput) error {
				return output.Decode(&v)
			}); err != nil {
				logger.Fatalf("number.folat ==> %s", err.Error())
			}

			logger.Debugf(`number.folat ==> %f`, v)
		}
	}

	// boolean
	{
		if hook, err := cli.GetObject(mxdata, key("boolean")); err != nil {
			logger.Fatalf("boolean ==> %s", err.Error())
		} else {
			var v bool

			if err := hook.Fetch(func(output storage.ObjectOutput) error {
				return output.Decode(&v)
			}); err != nil {
				logger.Fatalf("boolean ==> %s", err.Error())
			}

			logger.Debugf(`boolean ==> %v`, v)
		}
	}

	// json
	{
		if hook, err := cli.GetObject(mxdata, key("object.json")); err != nil {
			logger.Fatalf("object.json ==> %s", err.Error())
		} else {
			var v map[string]string

			if err := hook.Fetch(func(output storage.ObjectOutput) error {
				return output.Decode(&v)
			}); err != nil {
				logger.Fatalf("object.json ==> %s", err.Error())
			}

			logger.Debugf(`object.json ==> %v`, v)
		}
	}
}

func TestS3Client_DelObject(t *testing.T) {
	cli, _ := clientBOS()

	// put
	if err := cli.PutObject(storage.InputNumber(mxdata, key("del.int"), time.Now().Second())); err != nil {
		logger.Fatalf("del.int ==> %s", err.Error())
	}
	if found, err := cli.IsExist(mxdata, key("del.int")); err != nil {
		logger.Fatalf("del.int ==> %s", err.Error())
	} else {
		logger.Debugf("del.int ==> %v", found)
	}

	// del
	if err := cli.DelObject(mxdata, key("del.int")); err != nil {
		logger.Fatalf("del.int ==> %s", err.Error())
	}
	if found, err := cli.IsExist(mxdata, key("del.int")); err != nil {
		logger.Fatalf("del.int ==> %s", err.Error())
	} else {
		logger.Debugf("del.int ==> %v", found)
	}
}

func TestS3Client_DelPrefix(t *testing.T) {
	cli, _ := clientBOS()

	if err := cli.DelPrefix(mxdata, key("folder")); err != nil {
		logger.Fatal(err)
	}
}

func TestS3Client_PutFolder(t *testing.T) {
	cli, _ := clientBOS()

	if err := cli.PutFolder(mxdata, key("20221221-011/"), `C:\Users\sun\Desktop\20221221-011`); err != nil {
		logger.Fatal(err)
	}
}

func TestS3Client_GetObjectsWithIterator(t *testing.T) {
	cli, _ := clientBOS()

	var count int
	if err := cli.GetObjectsWithIterator(mxdata, key("20221221-011"), func(keys []*string) error {
		count += len(keys)

		for _, key := range keys {
			fmt.Println(*key)
		}
		return nil
	}); err != nil {
		logger.Fatal(err)
	}

	fmt.Println("==>", count)
}

func TestS3Client_Copy(t *testing.T) {
	cli, _ := clientBOS()

	err := cli.Copy(storage.PositionRemote(mxdata, key("simple/")), storage.PositionRemote(mxdata, key("simple1/")))
	if err != nil {
		logger.Fatal(err)
	}
}

func TestS3Client_Compress(t *testing.T) {
	cli, _ := clientBOS()

	f, e := os.Create(`data.tar.gz`)
	if e != nil {
		logger.Fatal(e)
	}
	defer f.Close()

	if err := cli.Compress(mxdata, key("20221221-011/"), f); err != nil {
		logger.Fatal(err)
	}
}
