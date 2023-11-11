package storage

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"go.uber.org/zap"
	"io"
)

type (
	Storage struct {
		bucketName        *string
		clientS3          *s3.Client
		clientS3Presigned *s3.PresignClient
		Logger            *zap.Logger
	}

	GetSignaturesOutput struct {
		KeyId          *string
		PublicKeyArmor *string
	}

	GetObjectOutput struct {
		Body   *io.ReadCloser
		Length int64
	}

	Metadata struct {
		KeyId  string
		Sha256 string
	}
)
