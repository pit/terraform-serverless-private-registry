package helpers

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ResponseNotFound struct {
	Status    string `json:"status"`
	Details   string `json:"details"`
	RequestId string `json:"requestId"`
}

type ResponseConflict struct {
	Status    string `json:"status"`
	Details   string `json:"details"`
	RequestId string `json:"requestId"`
}

func ApiErrorNotFound(requestId string, details string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusNotFound

	status := "Not found"
	stringBody, err := json.Marshal(&ResponseNotFound{
		Status:    status,
		Details:   details,
		RequestId: requestId,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func ApiErrorConflict(requestId string, details string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusConflict

	status := "Already exists"
	stringBody, err := json.Marshal(&ResponseConflict{
		Status:    status,
		Details:   details,
		RequestId: requestId,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func ApiErrorBadRequest(requestId string, details string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusBadRequest

	status := "Wrong request"
	stringBody, err := json.Marshal(&ResponseConflict{
		Status:    status,
		Details:   details,
		RequestId: requestId,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func ApiErrorIncorrectVersion(requestId string, details string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusNotAcceptable

	status := "Version format incorrect"
	stringBody, err := json.Marshal(&ResponseConflict{
		Status:    status,
		Details:   details,
		RequestId: requestId,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func ApiErrorUnknown(requestId string) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusInternalServerError

	status := "Unknown"
	stringBody, err := json.Marshal(&ResponseNotFound{
		Status:    status,
		RequestId: requestId,
	})
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response json: %s", string(stringBody))

	return &resp
}

func ApiErrorNoContent() *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = http.StatusNoContent

	return &resp
}

func ApiResponse(status int, body interface{}) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = status

	stringBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	resp.Body = string(stringBody)

	fmt.Printf("response: status=%d, json: %s", resp.StatusCode, string(stringBody))

	return &resp
}

func ApiResponseNoContent(status int) *events.APIGatewayProxyResponse {
	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{"Content-Type": "application/json"},
	}
	resp.StatusCode = status

	fmt.Printf("response: status=%d", resp.StatusCode)

	return &resp
}

func InitLogger(logLevel string, structured bool) (*zap.Logger, error) {
	spew.Config.Indent = "  "
	spew.Config.DisableMethods = true
	spew.Config.DisablePointerMethods = true

	var lvl zap.AtomicLevel
	switch strings.ToLower(logLevel) {
	case "info":
		lvl = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn", "warning":
		lvl = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error", "err":
		lvl = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		lvl = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	default:
		lvl = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	encoding := "console"
	if structured {
		encoding = "json"
	}

	cfg := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         encoding,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()

	return logger, err

}

func GenerateMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func GenerateMD5List(texts []string) []string {
	result := make([]string, len(texts))
	for idx, val := range texts {
		result[idx] = GenerateMD5(val)
	}
	return result
}

func ReaderToString(reader *io.ReadCloser) (*string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, *reader)
	if err != nil {
		return nil, errors.Wrap(err, "Error reading io.Reader to string")
	}
	result := buf.String()
	return &result, nil
}

func IntegrationTestSetup() {
	baseDir := os.Getenv("BASE_DIR")
	bucketName := os.Getenv("BUCKET_NAME")
	cmd := exec.Command("python3", fmt.Sprintf("%s/integration-test/init.py", baseDir), "--bucket", bucketName)
	stdOut, err := cmd.CombinedOutput()

	if stdOut != nil {
		fmt.Print(string(stdOut))
	}
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}

	fmt.Printf("IntegrationTestSetup() finished\n")
}

func FindSha256Sum(sha256Sums *string, filename *string) (*string, error) {
	sha256SumsList := strings.Split(strings.TrimRight(*sha256Sums, "\n"), "\n")

	for _, line := range sha256SumsList {
		if strings.HasSuffix(line, *filename) {
			sha256 := strings.Fields(line)
			if len(sha256) != 2 {
				return nil, errors.Errorf("Wrong sha256 format in line '%s'", line)
			}
			return &sha256[0], nil
		}
	}

	return nil, errors.Errorf("Line for %[1]s not found", filename)
}

func DecodeBase64(data string) ([]byte, error) {
	result, err := base64.StdEncoding.DecodeString(data)
	return result, err
}

type LambdaType int

const (
	LambdaTypeAuthorizer LambdaType = iota
	LambdaTypeCustomModulesUpload
	LambdaTypeCustomProvidersChecksumsUpload
	LambdaTypeCustomProvidersUpload
	LambdaTypeDefault
	LambdaTypeDiscovery
	LambdaTypeIndex
	LambdaTypeModulesDownload
	LambdaTypeModulesLatestVersion
	LambdaTypeModulesList
	LambdaTypeModulesSearch
	LambdaTypeModulesVersions
	LambdaTypeProvidersDownload
	LambdaTypeProvidersVersions
)

var (
	lambdaTypes = map[string]LambdaType{
		"default": LambdaTypeDefault,
		"index":   LambdaTypeIndex,

		"authorizer": LambdaTypeAuthorizer,
		"discovery":  LambdaTypeDiscovery,

		"modules-download":       LambdaTypeModulesDownload,
		"modules-latest-version": LambdaTypeModulesLatestVersion,
		"modules-list":           LambdaTypeModulesList,
		"modules-search":         LambdaTypeModulesSearch,
		"modules-versions":       LambdaTypeModulesVersions,

		"providers-download": LambdaTypeProvidersDownload,
		"providers-versions": LambdaTypeProvidersVersions,

		"custom-modules-upload":             LambdaTypeCustomModulesUpload,
		"custom-providers-checksums-upload": LambdaTypeCustomProvidersChecksumsUpload,
		"custom-providers-upload":           LambdaTypeCustomProvidersUpload,
	}
)

func GetLambdaType() LambdaType {
	str := os.Getenv("LAMBDA_TYPE")
	c, ok := lambdaTypes[strings.ToLower(str)]
	if !ok {
		_, err := os.Stderr.WriteString(fmt.Sprintf("Unknown lambda type: %s", str))
		if err != nil {
			os.Exit(-1)
		} else {
			os.Exit(-2)
		}
	}
	return c
}
