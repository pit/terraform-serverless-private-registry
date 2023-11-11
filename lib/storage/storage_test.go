package storage

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"os"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	"testing"
)

func TestMain(m *testing.M) {
	helpers.IntegrationTestSetup()
	code := m.Run()
	os.Exit(code)
}

func TestNewStorage(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, err := NewStorage(bucketName, logger)

	if storage == nil {
		t.Error("storageSvc is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestListDirs(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	resp, err := storage.ListDirs("0000", "modules/")

	if err != nil {
		t.Error(err)
	}
	if len(*resp) == 0 {
		t.Error("Response is empty")
	}

	ok := false
	for _, item := range *resp {
		if item == "modules/test-namespace/" {
			ok = true
		}
	}
	if !ok {
		t.Error("Wrong resp content, prefix 'modules/test-namespace/' wasn't found")
	}
}
func TestListFiles(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	resp, err := storage.ListFiles("0000", "modules/")

	if err != nil {
		t.Error(err)
	}
	if len(*resp) == 0 {
		t.Error("Response is empty")
	}

	ok := false
	for _, item := range *resp {
		if item == "modules/test-namespace/test-name/test-provider/2.0.0/test-namespace-test-name-test-provider-2.0.0.tar.gz" {
			ok = true
		}
	}
	if !ok {
		t.Error("Wrong resp content, file 'modules/test-namespace/test-name/test-provider/2.0.0/test-namespace-test-name-test-provider-2.0.0.tar.gz' wasn't found")
	}
}

func TestGetDownloadUrl(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)
	key := "modules/test-namespace/test-name/test-provider/2.0.0/test-namespace-test-name-test-provider-2.0.0.tar.gz"
	filename := "test-namespace-test-name-test-provider-2.0.0.tar.gz"
	resp, err := storage.GetDownloadUrl("0000", key, filename)

	if err != nil {
		t.Error(err)
	}
	if resp == nil {
		t.Error("Response is nil")
	}

	httpResp, httpErr := http.Get(*resp)
	if httpErr != nil {
		t.Error(httpErr)
	}
	if httpResp.StatusCode != 200 {
		t.Errorf("Wrong response status %d", httpResp.StatusCode)
	}
}

func TestGetObject(t *testing.T) {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)
	key := "providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_SHA256SUMS.sig.pub"
	resp, err := storage.GetObject("0000", key)

	fmt.Printf(">>> %x", err)

	if err != nil {
		//err, ok := errors.Cause(*err).(stackTracer)
		//if !ok {
		//	panic("oops, err does not implement stackTracer")
		//}
		//st := err.StackTrace()

		//fmt.Printf(") // top two frames
		t.Errorf("Error: %+v", errors.WithStack(err))
	} else if resp == nil {
		t.Error("Response is nil")
	} else {
		if resp.Body == nil {
			t.Error("Response body is nil")
		}

		if resp.Length == 0 {
			t.Error("Response has 0 length")
		}

		buf := new(strings.Builder)
		ioLength, ioErr := io.Copy(buf, *resp.Body)

		if ioErr != nil {
			t.Error("Error reading object body")
		}

		if ioLength == 0 {
			t.Error("Error reading object body - 0 bytes read")
		}
	}
}

func TestGetDownloadUrlNotFound(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	key := "not_exist"
	filename := "unknown"
	resp, err := storage.GetDownloadUrl("0000", key, filename)

	if err != nil {
		t.Error("Error is nil")
	}
	if resp == nil {
		t.Error("Response is not empty")
	}
}

func TestCheckObjectExist(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := NewStorage(bucketName, logger)

	key := "not_exist"
	resp, err := storage.CheckObjectExist("0000", key)

	if err != nil {
		t.Error("Error is nil")
	}
	if resp == nil {
		t.Error("Response is not empty")
	}
	if *resp == true {
		t.Error("Wrong response = should be false")
	}
}
