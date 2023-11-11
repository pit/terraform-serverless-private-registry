package providers

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/storage"
	"testing"
)

func TestMain(m *testing.M) {
	helpers.IntegrationTestSetup()
	code := m.Run()
	os.Exit(code)
}

func TestProvidersNewProviders(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage.NewStorage(bucketName, logger)
	modules, err := NewProviders(storage, logger)

	if modules == nil {
		t.Error("Providers is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}
}

func TestProvidersListProviderVersions(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage.NewStorage(bucketName, logger)
	modules, _ := NewProviders(storage, logger)

	providerNamespace := "test-namespace"
	providerType := "test-type"
	params := ListProviderVersionsInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
	}
	resp, err := modules.ListProviderVersions("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if resp.Versions == nil {
		t.Error("No modules in response")
	}

	ok := false
	for _, provider := range resp.Versions {
		if provider.Version == "" {
			t.Error("No versions in response.Providers")
		}

		for _, version := range resp.Versions {
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

func TestProvidersSaveChecksums(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage.NewStorage(bucketName, logger)
	providers, _ := NewProviders(storage, logger)

	providerNamespace := "test-namespace"
	providerType := "test-type"
	providerVersion := "2.1.1"

	strKeyId := readFixtureStr("provider-keyid")
	strSha256Sums := readFixtureStr("provider-sha256sums")
	strSha256SumsSig := readFixtureStrBase64("provider-sha256sums.sig")
	strSha256SumsSigPub := readFixtureStr("provider-sha256sums.sig.pub")

	params := SaveSignaturesInput{
		Namespace:        &providerNamespace,
		Type:             &providerType,
		Version:          &providerVersion,
		KeyId:            &strKeyId,
		Sha256Sums:       &strSha256Sums,
		Sha256SumsSig:    &strSha256SumsSig,
		Sha256SumsSigPub: &strSha256SumsSigPub,
	}

	resp, err := providers.SaveSignatures("0000", params)

	if err != nil {
		t.Error("Should not be error")
	}

	if resp == nil {
		t.Error("Response can't be nil")
	}
}

func TestProvidersGetDownload(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage.NewStorage(bucketName, logger)
	providers, _ := NewProviders(storage, logger)

	providerNamespace := "test-namespace"
	providerType := "test-type"
	providerVersion := "2.0.1"
	providerOS := "darwin"
	providerArch := "amd64"
	params := GetDownloadInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
		Version:   &providerVersion,
		OS:        &providerOS,
		Arch:      &providerArch,
	}
	resp, err := providers.GetDownload("0000", params)

	if resp == nil {
		t.Error("Response is nil")
	}
	if err != nil {
		t.Error("Error", err)
	}

	if !strings.HasPrefix(resp.DownloadUrl, fmt.Sprintf("https://%s.s3.eu-central-1.amazonaws.com/providers/test-namespace/test-type/2.0.1/terraform-provider-test-type_2.0.1_darwin_amd64.zip?", bucketName)) {
		t.Error("Wrong url format")
	}
}

func TestProvidersGetDownloadUrl404(t *testing.T) {
	bucketName := os.Getenv("BUCKET_NAME")
	logger, _ := helpers.InitLogger("DEBUG", true)
	storage, _ := storage.NewStorage(bucketName, logger)
	providers, _ := NewProviders(storage, logger)

	providerNamespace := "test-namespace"
	providerType := "test-type"
	providerVersion := "0.0.0"
	providerOS := "darwin"
	providerArch := "amd64"

	params := GetDownloadInput{
		Namespace: &providerNamespace,
		Type:      &providerType,
		Version:   &providerVersion,
		OS:        &providerOS,
		Arch:      &providerArch,
	}
	resp, err := providers.GetDownload("0000", params)

	if resp != nil {
		t.Error("Response is not nil")
	}
	if err == nil {
		t.Error("Should be error")
	}
}

func readFixtureStr(path string) string {
	data, err := os.ReadFile(fmt.Sprintf("%s/integration-test/fixtures/%s", os.Getenv("BASE_DIR"), path))
	if err != nil {
		panic(err)
	}
	result := string(data)
	return result
}

func readFixtureStrBase64(path string) string {
	data, err := os.ReadFile(fmt.Sprintf("%s/integration-test/fixtures/%s", os.Getenv("BASE_DIR"), path))
	if err != nil {
		panic(err)
	}
	result := base64.StdEncoding.EncodeToString(data)
	return result
}
