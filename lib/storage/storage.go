package storage

import (
	"bytes"
	"context"
	"fmt"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	//"github.com/aws/smithy-go/transport/http"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

//const (
//	ErrUnknown = iota
//	ErrObjectNotFound
//	ErrObjectNotAccessible
//)

//type StorageError struct {
//	Message    string
//	Code       int
//	BucketName string
//	Key        string
//	Err        error
//}
//
//func (s StorageError) Error() string {
//	panic(s.Message)
//}

const PresignUrlDuration = time.Duration(10) * time.Hour

func NewStorage(bucketName string, logger *zap.Logger) (*Storage, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		logger.Error("error loading aws config", zap.Error(err))
		return nil, fmt.Errorf("error loading aws config")
	}

	clientS3 := s3.NewFromConfig(awsCfg)
	clientS3Presigned := s3.NewPresignClient(clientS3)

	return &Storage{
		clientS3:          clientS3,
		clientS3Presigned: clientS3Presigned,
		bucketName:        &bucketName,
		Logger:            logger,
	}, nil
}

func (svc *Storage) ListDirs(ctxId string, key string) (*[]string, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.ListDirs() called", ctxId),
		zap.String("key", key),
	)

	if !strings.HasSuffix(key, "/") {
		key = fmt.Sprintf("%s/", key)
	}

	delimiter := "/"
	params := &s3.ListObjectsV2Input{
		Bucket:    svc.bucketName,
		Prefix:    &key,
		Delimiter: &delimiter,
	}
	resp, err := svc.clientS3.ListObjectsV2(context.Background(), params)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}
	svc.Logger.Debug(fmt.Sprintf("%s storage.ListFiles().ListObjectsV2()", ctxId),
		zap.Reflect("params", params),
		zap.Reflect("resp", resp),
	)

	var result []string
	for _, path := range resp.CommonPrefixes {
		result = append(result, *path.Prefix)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.ListDirs() return", ctxId),
		zap.String("key", key),
		zap.Reflect("result", result),
	)

	return &result, nil
}

func (svc *Storage) ListFiles(ctxId string, key string) (*[]string, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.ListFiles() called", ctxId),
		zap.String("key", key),
	)

	if !strings.HasSuffix(key, "/") {
		key = fmt.Sprintf("%s/", key)
	}

	//delimiter := "/"
	params := &s3.ListObjectsV2Input{
		Bucket: svc.bucketName,
		Prefix: &key,
		//Delimiter: &delimiter,
	}
	resp, err := svc.clientS3.ListObjectsV2(context.Background(), params)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}
	svc.Logger.Debug(fmt.Sprintf("%s storage.ListFiles().ListObjectsV2()", ctxId),
		zap.Reflect("params", params),
		zap.Reflect("resp", resp),
	)

	var result []string

	for _, path := range resp.Contents {
		result = append(result, *path.Key)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.ListFiles() return", ctxId),
		zap.String("key", key),
		zap.Reflect("result", result),
	)

	return &result, nil
}

//func (svc *Storage) GetMetadata(ctxId string, base string, provider_type string, provider_version string) (*Metadata, error) {
//	svc.Logger.Debug(fmt.Sprintf("%s storage.GetMetadata() called", ctxId),
//		zap.String("base", base),
//		zap.String("file", file),
//	)
//
//	key := fmt.Sprintf("%[1]s/keyid.txt", base, provider_type, provider_version)
//	params := &s3.GetObjectInput{
//		Bucket: svc.bucketName,
//		Key:    &key,
//	}
//	resp, err := svc.clientS3.GetObject(context.Background(), params)
//	if err != nil {
//		return nil, svc.handleError(ctxId, err)
//	}
//	//result := resp.Metadata
//
//	svc.Logger.Debug(fmt.Sprintf("%s storage.GetMetadata() return", ctxId),
//		zap.String("key", key),
//		zap.Reflect("result", result),
//	)
//
//	result := Metadata{}
//
//	return &result, nil
//}
//
//func (svc *Storage) SaveMetadata(ctxId string, key string, metadata *map[string]string) error {
//	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveMetadata() called", ctxId),
//		zap.String("key", key),
//		zap.Reflect("metadata", metadata),
//	)
//
//	headParams := &s3.HeadObjectInput{
//		Bucket: svc.bucketName,
//		Key:    &key,
//	}
//	headResp, headErr := svc.clientS3.HeadObject(context.Background(), headParams)
//	if headErr != nil {
//		return svc.handleError(ctxId, headErr)
//	}
//	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveMetadata().HeadObject", ctxId),
//		zap.Reflect("resp", headResp),
//	)
//
//	if headResp.Metadata == nil {
//		headResp.Metadata = map[string]string{}
//	}
//
//	for k, v := range *metadata {
//		headResp.Metadata[k] = v
//	}
//
//	copySource := fmt.Sprintf("%s/%s", *svc.bucketName, key)
//	copyParams := &s3.CopyObjectInput{
//		Bucket:     svc.bucketName,
//		Key:        &key,
//		CopySource: &copySource,
//		Metadata:   headResp.Metadata,
//	}
//
//	copyResp, copyErr := svc.clientS3.CopyObject(context.Background(), copyParams)
//	if copyErr != nil {
//		return svc.handleError(ctxId, copyErr)
//	}
//	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveMetadata().CopyObject", ctxId),
//		zap.Reflect("copyResp", copyResp),
//	)
//
//	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveMetadata() return", ctxId),
//		zap.String("key", key),
//	)
//	return nil
//}

func (svc *Storage) GetDownloadUrl(ctxId string, key string, fileName string) (*string, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.GetDownloadUrl() called", ctxId),
		zap.String("key", key),
		zap.String("fileName", fileName),
	)

	//paramsHead := &s3.HeadObjectInput{
	//	Bucket: svc.bucketName,
	//	Key:    &key,
	//}
	//_, err := svc.clientS3.HeadObject(context.TODO(), paramsHead)
	//if err != nil {
	//	return nil, svc.handleError(ctxId, err, "storageSvc.GetDownloadUrl.HeadObject", key,
	//		zap.Reflect("bucket", svc.bucketName),
	//		zap.String("key", key),
	//		zap.Reflect("params", paramsHead),
	//	)
	//}

	contentDisposition := fmt.Sprintf("attachment; filename=%s", fileName)
	contentType := "application/x-gzip"
	cacheControl := fmt.Sprintf("max-age=%.0f", PresignUrlDuration.Seconds())
	paramsSign := &s3.GetObjectInput{
		Bucket:                     svc.bucketName,
		Key:                        &key,
		ResponseContentDisposition: &contentDisposition,
		ResponseCacheControl:       &cacheControl,
		ResponseContentType:        &contentType,
	}

	resp, err := svc.clientS3Presigned.PresignGetObject(context.TODO(), paramsSign, s3.WithPresignExpires(PresignUrlDuration))
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.GetDownloadUrl() return", ctxId),
		zap.String("key", key),
		zap.String("fileName", fileName),
		zap.Reflect("result", &resp.URL),
	)

	return &resp.URL, nil
}

func (svc *Storage) GetUploadUrl(ctxId string, key string, fileName string) (*string, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.GetUploadUrl() called", ctxId),
		zap.String("key", key),
	)

	//contentDisposition := fmt.Sprintf("attachment; filename=%s", fileName)
	//contentType := "application/x-gzip"
	//cacheControl := fmt.Sprintf("max-age=%.0f", PresignUrlDuration.Seconds())
	paramsSign := &s3.PutObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
	}

	resp, err := svc.clientS3Presigned.PresignPutObject(context.TODO(), paramsSign, s3.WithPresignExpires(PresignUrlDuration))
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.GetUploadUrl() return", ctxId),
		zap.String("key", key),
		zap.Reflect("result", &resp.URL),
	)

	return &resp.URL, nil
}

var (
	BoolTrue  = true
	BoolFalse = false
)

func (svc *Storage) CheckObjectExist(ctxId string, key string) (*bool, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.CheckObjectExist() called", ctxId),
		zap.String("key", key),
	)

	paramsHead := &s3.HeadObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
	}

	respHeadObject, err := svc.clientS3.HeadObject(context.TODO(), paramsHead)
	if err != nil {
		var responseError *awshttp.ResponseError
		if errors.As(err, &responseError) && responseError.ResponseError.HTTPStatusCode() == http.StatusNotFound {
			svc.Logger.Debug(fmt.Sprintf("%s storage.CheckObjectExist() return false", ctxId),
				zap.String("key", key),
				zap.Reflect("params", paramsHead),
				zap.Error(err),
			)

			return &BoolFalse, nil
		} else {
			return nil, svc.handleError(ctxId, err)
		}
	}
	if respHeadObject != nil {
		svc.Logger.Debug(fmt.Sprintf("%s storage.CheckObjectExist() return true", ctxId),
			zap.String("key", key),
			zap.Reflect("params", paramsHead),
			zap.Reflect("resp", respHeadObject),
		)

		return &BoolTrue, nil
	}

	return nil, svc.handleError(ctxId, err)
}

func (svc *Storage) GetObject(ctxId string, key string) (*GetObjectOutput, error) {
	svc.Logger.Debug(fmt.Sprintf("%s storage.GetObject() called", ctxId),
		zap.String("key", key),
	)

	//paramsHead := &s3.HeadObjectInput{
	//	Bucket: svc.bucketName,
	//	Key:    &key,
	//}
	//_, err := svc.clientS3.HeadObject(context.TODO(), paramsHead)
	//if err != nil {
	//	return nil, svc.handleError(ctxId, err)
	//}

	paramsGet := &s3.GetObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
	}

	resp, err := svc.clientS3.GetObject(context.TODO(), paramsGet)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}
	svc.Logger.Debug(fmt.Sprintf("%s storage.GetObject().GetObject", ctxId),
		zap.String("key", key),
		zap.Int64("body.ContentLength", resp.ContentLength),
		zap.Reflect("resp", resp),
	)

	//buf := make([]byte, resp.ContentLength)
	//var readLen int
	//readLen, err = resp.Body.Read(buf)
	//if err != nil {
	//	svc.Logger.Debug(fmt.Sprintf("%s Reading object body", ctxId),
	//		zap.Int64("body.ContentLength", resp.ContentLength),
	//		zap.Int("readLen", readLen),
	//	)
	//	return nil, svc.handleError(ctxId, err)
	//}

	result := GetObjectOutput{
		Body:   &resp.Body,
		Length: resp.ContentLength,
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.GetObject() return", ctxId),
		zap.String("key", key),
		zap.Int64("contentLength", resp.ContentLength),
		zap.Reflect("result", result),
	)

	return &result, nil
}

func (svc *Storage) SaveObject(ctxId string, key string, data []byte) error {
	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveObject() called", ctxId),
		zap.String("key", key),
	)

	reader := bytes.NewReader(data)
	paramsPut := &s3.PutObjectInput{
		Bucket: svc.bucketName,
		Key:    &key,
		Body:   reader,
	}

	_, err := svc.clientS3.PutObject(context.TODO(), paramsPut)
	if err != nil {
		return svc.handleError(ctxId, err)
	}

	svc.Logger.Debug(fmt.Sprintf("%s storage.SaveObject() return", ctxId),
		zap.String("key", key),
	)

	return nil
}

func (svc *Storage) handleError(ctxId string, err error) error {
	// We print debug message here.
	// Since Error message should be printed in the final point - when whole call stack was started
	//svc.Logger.Debug(fmt.Sprintf("%s storage error", ctxId),
	//	zap.Error(err),
	//	zap.String("ctxId", ctxId),
	//)

	wrappedError := errors.Wrapf(err, "%s storage error", ctxId)

	return wrappedError
}
