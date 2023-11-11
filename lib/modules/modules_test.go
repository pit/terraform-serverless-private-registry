package modules

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	storage2 "terraform-serverless-private-registry/lib/storage"
	"testing"
)

func TestMain(m *testing.M) {
	helpers.IntegrationTestSetup()
	code := m.Run()
	os.Exit(code)
}

func TestModulesNewModules(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, err := NewModules(storage, logger)

	if modules == nil {
		t.Error("Modules is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestModulesListModuleVersions(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "test-namespace"
	name := "test-name"
	provider := "test-provider"
	params := ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
	}
	resp, err := modules.ListModuleVersions("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if resp.ModuleExists == false {
		t.Error("Response ModuleExists should be true")
	}

	if resp.Modules == nil {
		t.Error("No modules in response")
	}

	ok := false
	for _, module := range resp.Modules {
		if module.Versions == nil {
			t.Error("No versions in response.Modules")
		}

		for _, version := range module.Versions {
			m, err := regexp.MatchString(`\d+\.\d+\.\d+`, version.Version)
			if err != nil {
				t.Error(err)
			}
			if m {
				ok = true
			}
		}
	}
	if !ok {
		t.Error("Wrong resp content")
	}

}

func TestModulesGetDownloadUrl404(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "test-namespace"
	name := "test-name"
	provider := "test-provider"
	version := "2.0.0"
	params := ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, err := modules.GetDownloadUrl("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error is not nil", err)
	}
	if resp.ModuleExists != true {
		t.Error("Response ModuleExists shouldn't be true")
	}
}

func TestModulesGetDownloadUrl(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage2.NewStorage(bucketName, logger)
	modules, _ := NewModules(storage, logger)

	namespace := "test-namespace"
	name := "test-name"
	provider := "test-provider"
	version := "0.0.0"
	params := ModuleParams{
		Namespace: &namespace,
		Name:      &name,
		Provider:  &provider,
		Version:   &version,
	}
	resp, err := modules.GetDownloadUrl("0000", params)

	if err != nil {
		t.Error("Should be error")
	}
	if resp == nil {
		t.Error("Response is not nil")
	}

	if resp.Url == nil {
		t.Error("Response Url is nil")
	} else if !strings.HasPrefix(*resp.Url, fmt.Sprintf("https://%s.s3.eu-central-1.amazonaws.com/modules/test-namespace/test-name/test-provider/0.0.0/test-namespace-test-name-test-provider-0.0.0.tar.gz?", bucketName)) {
		t.Error("Wrong url format")
	}
}
