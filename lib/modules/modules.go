package modules

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
	"terraform-serverless-private-registry/lib/storage"
)

type Modules struct {
	storageSvc *storage.Storage
	logger     *zap.Logger
}

func NewModules(storage *storage.Storage, log *zap.Logger) (*Modules, error) {
	return &Modules{
		storageSvc: storage,
		logger:     log,
	}, nil
}

func (svc *Modules) ListModuleVersions(ctxId string, params ModuleParams) (*ListModuleVersionsResponse, error) {
	svc.logger.Debug(fmt.Sprintf("%s modules.ListModuleVersions() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
	)
	dirPath := fmt.Sprintf("modules/%s/%s/%s/", *params.Namespace, *params.Name, *params.Provider)
	dirs, err := svc.storageSvc.ListDirs(ctxId, dirPath)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}
	svc.logger.Debug(fmt.Sprintf("%s modules.ListModuleVersions().ListDirs() called", ctxId),
		zap.Strings("provider", *dirs),
	)

	var versions []Version
	for _, dir := range *dirs {
		dir = strings.TrimPrefix(dir, dirPath)
		dir = strings.TrimSuffix(dir, "/")
		versions = append(versions, Version{Version: dir})
	}

	result := ListModuleVersionsResponse{
		ModuleExists: true,
		Modules:      []ModuleVersions{{Versions: versions}},
	}

	svc.logger.Debug(fmt.Sprintf("%s ListModuleVersions() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Modules) GetDownloadUrl(ctxId string, params ModuleParams) (*GetDownloadUrlOutput, error) {
	svc.logger.Debug(fmt.Sprintf("%s modules.GetDownloadUrl() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
		zap.String("version", *params.Version),
	)
	key := fmt.Sprintf("modules/%[1]s/%[2]s/%[3]s/%[4]s/%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)
	fileName := fmt.Sprintf("%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)

	objectExist, objectExistErr := svc.storageSvc.CheckObjectExist(ctxId, key)
	if objectExistErr != nil {
		return nil, svc.handleError(ctxId, objectExistErr)
	}

	var result GetDownloadUrlOutput
	if *objectExist == true {
		url, err := svc.storageSvc.GetDownloadUrl(ctxId, key, fileName)
		if err != nil {
			return nil, svc.handleError(ctxId, err)
		}

		result = GetDownloadUrlOutput{
			ModuleExists: true,
			Url:          url,
		}
	} else {
		result = GetDownloadUrlOutput{
			ModuleExists: false,
		}
	}

	svc.logger.Debug(fmt.Sprintf("%s GetDownloadUrl() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Modules) GetUploadUrl(ctxId string, params ModuleParams) (*UploadModuleResponse, error) {
	svc.logger.Debug(fmt.Sprintf("%s modules.GetUploadUrl() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Name),
		zap.String("provider", *params.Provider),
		zap.String("version", *params.Version),
	)

	key := fmt.Sprintf("modules/%[1]s/%[2]s/%[3]s/%[4]s/%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)
	fileName := fmt.Sprintf("%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)

	objectExist, objectExistErr := svc.storageSvc.CheckObjectExist(ctxId, key)
	if objectExistErr != nil {
		return nil, objectExistErr
	}

	var result UploadModuleResponse
	if *objectExist {
		result = UploadModuleResponse{
			ModuleExists: true,
		}
	} else {
		resp, err := svc.storageSvc.GetUploadUrl(ctxId, key, fileName)
		if err != nil {
			return nil, svc.handleError(ctxId, err)
		}

		result = UploadModuleResponse{
			Url: resp,
		}
	}

	svc.logger.Debug(fmt.Sprintf("%s GetUploadUrl() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

//func (svc *Modules) CheckModuleExist(ctxId string, params ModuleParams) (*bool, error) {
//	svc.logger.Debug(fmt.Sprintf("%s modules.GetUploadUrl() called", ctxId),
//		zap.String("namespace", *params.Namespace),
//		zap.String("name", *params.Name),
//		zap.String("provider", *params.Provider),
//		zap.String("version", *params.Version),
//	)
//
//	key := fmt.Sprintf("modules/%[1]s/%[2]s/%[3]s/%[4]s/%[1]s-%[2]s-%[3]s-%[4]s.tar.gz", *params.Namespace, *params.Name, *params.Provider, *params.Version)
//	keyExists, keyExistError := svc.storageSvc.CheckObjectExist(ctxId, key)
//
//	if keyExistError != nil {
//		svc.logger.Debug(fmt.Sprintf("%s modules.GetUploadUrl() error checking key in bucket", ctxId),
//			zap.String("key", key),
//		)
//
//		return nil, svc.handleError(ctxId, keyExistError, "GetUploadUrl", params)
//	}
//
//	if *keyExists {
//		svc.logger.Debug(fmt.Sprintf("%s modules.GetUploadUrl() key exist in bucket", ctxId),
//			zap.String("key", key),
//		)
//
//		return &BoolTrue, nil
//	}
//
//	return &BoolFalse, nil
//}

func (svc *Modules) handleError(ctxId string, err error) error {
	return errors.Wrapf(err, "%s modules error", ctxId)
}
