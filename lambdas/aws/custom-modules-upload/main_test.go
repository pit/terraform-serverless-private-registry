package aws_custome_modules_upload

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"io/ioutil"
	"net/url"
	"os"
	"path"
	"terraform-serverless-private-registry/lib/helpers"
	"testing"
)

func TestMain(m *testing.M) {
	Setup()
	helpers.IntegrationTestSetup()
	code := m.Run()
	os.Exit(code)
}

func readBaseRequest() events.APIGatewayProxyRequest {
	baseDir := os.Getenv("BASE_DIR")
	file, _ := ioutil.ReadFile(path.Join(baseDir, "integration-test/fixtures/api-gateway-request.json"))

	result := events.APIGatewayProxyRequest{}
	_ = json.Unmarshal([]byte(file), &result)

	return result
}

func parseResponse(data string) map[string]interface{} {
	var result map[string]interface{}
	_ = json.Unmarshal([]byte(data), &result)
	return result
}

func TestBase(t *testing.T) {
	for _, version := range []string{"1.2.1", "1.2.2", "1.2.3"} {
		request := readBaseRequest()
		request.PathParameters = map[string]string{
			"namespace": "test-namespace",
			"name":      "test-name",
			"provider":  "test-provider",
			"version":   version,
		}

		resp, err := Handler(context.TODO(), request)

		if err != nil {
			t.Errorf("got error: %q", err)
		}

		if resp == nil {
			t.Error("Response is nil")
		}

		if resp.StatusCode != 200 {
			t.Errorf("Status code - expected 200, got %d", resp.StatusCode)
		}

		respDict := parseResponse(resp.Body)

		if val, ok := respDict["url"]; ok {
			_, err := url.ParseRequestURI(fmt.Sprintf("%v", val))
			if err != nil {
				t.Errorf("Url is not valid")
			}
		} else {
			t.Errorf("Response doesn't have url field")
		}
	}
}

func TestLatest(t *testing.T) {
	version := "latest"

	request := readBaseRequest()
	request.PathParameters = map[string]string{
		"namespace": "test-namespace",
		"name":      "test-name",
		"provider":  "test-provider",
		"version":   version,
	}

	resp, err := Handler(context.TODO(), request)

	if err != nil {
		t.Errorf("got error: %q", err)
	}

	if resp == nil {
		t.Error("Response is nil")
	}

	if resp.StatusCode != 406 {
		t.Errorf("Status code - expected 406, got %d", resp.StatusCode)
	}

	respDict := parseResponse(resp.Body)

	if val, ok := respDict["url"]; ok {
		_, err := url.ParseRequestURI(fmt.Sprintf("%v", val))
		if err != nil {
			t.Errorf("Url is not valid")
		}
	} else {
		t.Errorf("Response doesn't have url field")
	}
}

func TestExist(t *testing.T) {
	request := readBaseRequest()
	request.PathParameters = map[string]string{
		"namespace": "test-namespace",
		"name":      "test-name",
		"provider":  "test-provider",
		"version":   "2.0.0",
	}

	resp, err := Handler(context.TODO(), request)

	if err != nil {
		t.Errorf("got error: %q", err)
	}

	if resp == nil {
		t.Error("Response is nil")
	}

	if resp.StatusCode != 409 {
		t.Errorf("Status code - expected 409, got %d", resp.StatusCode)
	}

	respDict := parseResponse(resp.Body)

	if _, ok := respDict["url"]; ok {
		t.Errorf("Response has url field, but shouldn't")
	}
}
