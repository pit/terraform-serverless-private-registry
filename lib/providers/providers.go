package providers

import (
	"encoding/base64"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"regexp"
	"strings"
	"terraform-serverless-private-registry/lib/helpers"
	"terraform-serverless-private-registry/lib/storage"
)

type Providers struct {
	storageSvc *storage.Storage
	logger     *zap.Logger
}

type ListProviderVersionsOutput struct {
	ProviderExists *bool
	Versions       []ProviderVersion `json:"versions"`
}

type ProviderVersion struct {
	Version   string             `json:"version"`
	Protocols []string           `json:"protocols"`
	Platforms []ProviderPlatform `json:"platforms"`
}
type ProviderPlatform struct {
	OS   string `json:"os"`
	Arch string `json:"arch"`
}

type ListProviderVersionsInput struct {
	Namespace *string
	Type      *string
}

type GetDownloadInput struct {
	Namespace *string
	Type      *string
	Version   *string
	OS        *string
	Arch      *string
}

type GpgPublicKey struct {
	KeyId      string `json:"key_id"`
	AsciiArmor string `json:"ascii_armor"`
	//TrustSignature string `json:"trust_signature"`
	Source    string `json:"source"`
	SourceUrl string `json:"source_url"`
}

type SigningKeys struct {
	GpgPublicKeys []GpgPublicKey `json:"gpg_public_keys"`
}

type GetDownloadOutput struct {
	ProviderExists *bool       `json:"-"`
	Protocols      []string    `json:"protocols"`
	OS             string      `json:"os"`
	Arch           string      `json:"arch"`
	Filename       string      `json:"filename"`
	DownloadUrl    string      `json:"download_url"`
	ShaSumsUrl     string      `json:"shasums_url"`
	ShaSumsSigUrl  string      `json:"shasums_signature_url"`
	ShaSum         string      `json:"shasum"`
	SigningKeys    SigningKeys `json:"signing_keys"`
}

type GetUploadInput struct {
	Namespace *string
	Type      *string
	Version   *string
	OS        *string
	Arch      *string
	Sha256    *string
}

type UploadProviderResponse struct {
	ProviderExists *bool
	Url            string `json:"url"`
}

const (
	MetadataSha256 = "sha256"
	MetadataKeyId  = "keyId"
)

var (
	TruePtr  = true
	FalsePtr = false
)

func NewProviders(storage *storage.Storage, log *zap.Logger) (*Providers, error) {
	return &Providers{
		storageSvc: storage,
		logger:     log,
	}, nil
}

var DEFAULT_PROTOCOLS = []string{"5.0"}

func (svc *Providers) ListProviderVersions(ctxId string, params ListProviderVersionsInput) (*ListProviderVersionsOutput, error) {
	svc.logger.Debug("providers.ListProviderVersions() called",
		zap.String("ctxId", ctxId),
		zap.Reflect("params", params),
	)
	dirPath := fmt.Sprintf("providers/%s/%s/", *params.Namespace, *params.Type)
	files, err := svc.storageSvc.ListFiles(ctxId, dirPath)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}

	var versions []*ProviderVersion
	var curVer *ProviderVersion
	providerRegex := regexp.MustCompile("providers/(?P<namespace>.*?)\\/(?P<type>.*?)\\/(?P<version>.*?)\\/terraform-provider-.*_(?P<os>.*?)_(?P<arch>.*?)\\.zip")
	for _, file := range *files {
		if matches := providerRegex.FindAllStringSubmatch(file, -1); matches != nil {
			svc.logger.Debug(fmt.Sprintf("%s D1", ctxId),
				zap.Reflect("matches", matches),
				zap.Reflect("versions", versions),
				zap.Int("idx.type", providerRegex.SubexpIndex("namespace")),
				zap.Int("idx.type", providerRegex.SubexpIndex("type")),
				zap.Int("idx.version", providerRegex.SubexpIndex("version")),
				zap.Int("idx.os", providerRegex.SubexpIndex("os")),
				zap.Int("idx.arch", providerRegex.SubexpIndex("arch")),
			)

			matchVersion := matches[0][providerRegex.SubexpIndex("version")]
			matchOs := matches[0][providerRegex.SubexpIndex("os")]
			matchArch := matches[0][providerRegex.SubexpIndex("arch")]

			file = strings.TrimPrefix(file, dirPath)
			file = strings.TrimSuffix(file, "/")
			if curVer == nil || curVer.Version != matchVersion {
				curVer = &ProviderVersion{
					Version:   matchVersion,
					Protocols: DEFAULT_PROTOCOLS,
					Platforms: []ProviderPlatform{
						{
							OS:   matchOs,
							Arch: matchArch,
						},
					},
				}
				versions = append(versions, curVer)
			} else {
				curVer.Platforms = append(curVer.Platforms,
					ProviderPlatform{
						OS:   matchOs,
						Arch: matchArch,
					},
				)
			}
		}
	}

	versionsData := make([]ProviderVersion, 0)
	for _, ver := range versions {
		versionsData = append(versionsData, *ver)
	}
	result := ListProviderVersionsOutput{
		Versions: versionsData,
	}

	svc.logger.Debug(fmt.Sprintf("%s ListProviderVersions() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Providers) GetDownload(ctxId string, params GetDownloadInput) (*GetDownloadOutput, error) {
	svc.logger.Debug(fmt.Sprintf("%s Providers.GetDownload() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("type", *params.Type),
		zap.String("version", *params.Version),
		zap.String("os", *params.OS),
		zap.String("arch", *params.Arch),
	)

	downloadFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_%[3]s_%[4]s.zip", *params.Type, *params.Version, *params.OS, *params.Arch)
	downloadUrlKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_%[4]s_%[5]s.zip", *params.Namespace, *params.Type, *params.Version, *params.OS, *params.Arch)
	downloadUrl, errStorage := svc.storageSvc.GetDownloadUrl(ctxId, downloadUrlKey, downloadFileName)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	}

	shaSumsFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_SHA256SUMS", *params.Type, *params.Version)
	shaSumsKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS", *params.Namespace, *params.Type, *params.Version)
	shaSumsUrl, errStorage := svc.storageSvc.GetDownloadUrl(ctxId, shaSumsKey, shaSumsFileName)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	}

	shaSumsSigFileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_SHA256SUMS.sig", *params.Type, *params.Version)
	shaSumsSigKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig", *params.Namespace, *params.Type, *params.Version)
	shaSumsSigUrl, errStorage := svc.storageSvc.GetDownloadUrl(ctxId, shaSumsSigKey, shaSumsSigFileName)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	}

	keyAsciiArmorKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig.pub", *params.Namespace, *params.Type, *params.Version)
	keyAsciiArmorData, errStorage := svc.storageSvc.GetObject(ctxId, keyAsciiArmorKey)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	}
	keyAsciiArmor, keyAsciiArmorErr := helpers.ReaderToString(keyAsciiArmorData.Body)
	if keyAsciiArmorErr != nil {
		return nil, svc.handleError(ctxId, keyAsciiArmorErr)
	}

	metadataKeyIdKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig.keyid", *params.Namespace, *params.Type, *params.Version)
	metadataKeyIdData, errStorage := svc.storageSvc.GetObject(ctxId, metadataKeyIdKey)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	} else {
		svc.logger.Debug("GetDownload.Metadata",
			zap.Reflect("metadataKeyIdKey", metadataKeyIdKey),
		)
	}
	metadataKeyId, metadataKeyIdErr := helpers.ReaderToString(metadataKeyIdData.Body)
	if metadataKeyIdErr != nil {
		return nil, svc.handleError(ctxId, metadataKeyIdErr)
	}

	metadataSha256SumsData, errStorage := svc.storageSvc.GetObject(ctxId, shaSumsKey)
	if errStorage != nil {
		return nil, svc.handleError(ctxId, errStorage)
	} else {
		svc.logger.Debug("GetDownload.Metadata",
			zap.Reflect("metadataSha256Key", shaSumsKey),
		)
	}
	metadataSha256Sums, metadataSha256Err := helpers.ReaderToString(metadataSha256SumsData.Body)
	if metadataSha256Err != nil {
		return nil, svc.handleError(ctxId, metadataSha256Err)
	}
	metadataSha256, metadataSha256Err := helpers.FindSha256Sum(metadataSha256Sums, &downloadFileName)
	if metadataSha256Err != nil {
		return nil, svc.handleError(ctxId, metadataSha256Err)
	}

	result := GetDownloadOutput{
		Protocols:     DEFAULT_PROTOCOLS,
		OS:            *params.OS,
		Arch:          *params.Arch,
		Filename:      downloadFileName,
		DownloadUrl:   *downloadUrl,
		ShaSumsUrl:    *shaSumsUrl,
		ShaSumsSigUrl: *shaSumsSigUrl,
		ShaSum:        *metadataSha256,
		SigningKeys: SigningKeys{
			GpgPublicKeys: []GpgPublicKey{
				{
					KeyId:      *metadataKeyId,
					AsciiArmor: *keyAsciiArmor,
					Source:     "Kvinta Gmbh",
					SourceUrl:  "https://www.kvinta.com/",
				},
			},
		},
	}

	svc.logger.Debug(fmt.Sprintf("%s ListModuleVersions() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Providers) GetUpload(ctxId string, params GetUploadInput) (*UploadProviderResponse, error) {
	svc.logger.Debug(fmt.Sprintf("%s provideres.GetUpload() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("name", *params.Type),
		zap.String("version", *params.Version),
		zap.String("os", *params.OS),
		zap.String("arch", *params.Arch),
	)

	fileName := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_%[3]s_%[4]s.zip", *params.Type, *params.Version, *params.OS, *params.Arch)
	key := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/%[4]s", *params.Namespace, *params.Type, *params.Version, fileName)

	keyExists, keyExistsErr := svc.storageSvc.CheckObjectExist(ctxId, key)
	if keyExistsErr != nil {
		return nil, svc.handleError(ctxId, keyExistsErr)
	}
	if *keyExists {
		svc.logger.Debug(fmt.Sprintf("Key %[1]s exists", key))
		result := UploadProviderResponse{
			ProviderExists: &TruePtr,
		}
		return &result, nil
	}

	resp, err := svc.storageSvc.GetUploadUrl(ctxId, key, fileName)
	if err != nil {
		return nil, svc.handleError(ctxId, err)
	}

	result := UploadProviderResponse{
		Url: *resp,
	}

	svc.logger.Debug(fmt.Sprintf("%s GetUpload() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Providers) SaveSignatures(ctxId string, params SaveSignaturesInput) (*SaveSignaturesOutput, error) {
	svc.logger.Debug(fmt.Sprintf("%s provideres.SaveSignatures() called", ctxId),
		zap.String("namespace", *params.Namespace),
		zap.String("type", *params.Type),
		zap.String("version", *params.Version),
		zap.Reflect("params", params),
	)

	keyPrefix := fmt.Sprintf("providers/%s/%s/%s", *params.Namespace, *params.Type, *params.Version)
	listResp, listErr := svc.storageSvc.ListFiles(ctxId, keyPrefix)
	if listErr != nil {
		return nil, svc.handleError(ctxId, listErr)
	}
	if len(*listResp) == 0 {
		result := SaveSignaturesOutput{
			ProviderExists: &FalsePtr,
		}

		svc.logger.Debug(fmt.Sprintf("%s SaveSignatures() return", ctxId),
			zap.Reflect("result", result),
		)
		return &result, nil
	}

	shaSumsCount := 0
	sha256Sums := strings.Split(strings.TrimRight(*params.Sha256Sums, "\n"), "\n")
	keyRegexp := regexp.MustCompile(".*?\\/terraform-provider-(?P<type>.*?)_(?P<version>.*?)_(?P<os>.*?)_(?P<arch>.*?)\\.zip")
	for _, key := range *listResp {
		svc.logger.Debug(fmt.Sprintf("%s D", ctxId),
			zap.String("key", key),
		)

		if matches := keyRegexp.FindAllStringSubmatch(key, -1); matches != nil {
			svc.logger.Debug(fmt.Sprintf("%s D1", ctxId),
				zap.Reflect("matches", matches),
				zap.Int("idx.type", keyRegexp.SubexpIndex("type")),
				zap.Int("idx.version", keyRegexp.SubexpIndex("version")),
				zap.Int("idx.os", keyRegexp.SubexpIndex("os")),
				zap.Int("idx.arch", keyRegexp.SubexpIndex("arch")),
			)

			matchedType := matches[0][keyRegexp.SubexpIndex("type")]
			matchedVersion := matches[0][keyRegexp.SubexpIndex("version")]
			matchedOs := matches[0][keyRegexp.SubexpIndex("os")]
			matchedArch := matches[0][keyRegexp.SubexpIndex("arch")]
			suffix := fmt.Sprintf("terraform-provider-%[1]s_%[2]s_%[3]s_%[4]s.zip", matchedType, matchedVersion, matchedOs, matchedArch)

			svc.logger.Debug("Scanning Sha256Sums",
				zap.String("suffix", suffix),
				zap.Reflect("sha256sums", sha256Sums),
			)
			for _, line := range sha256Sums {
				svc.logger.Debug("sha256 line",
					zap.String("line", line),
					zap.String("suffix", suffix),
					zap.Int("shaSumsCount", shaSumsCount),
				)
				if strings.HasSuffix(line, suffix) {
					shaSumsCount++
					svc.logger.Debug("sha256 line match",
						zap.String("line", line),
						zap.Int("shaSumsCount", shaSumsCount),
					)
					sha256 := strings.Fields(line)
					if len(sha256) != 2 {
						return nil, errors.Errorf("Wrong sha256 format in line '%s'", line)
					}
				}
			}
		}
	}

	if len(*listResp) != shaSumsCount {
		details := fmt.Sprintf("Number of sha256sum lines(%d) not equal to number of provider files(%d) in storage", shaSumsCount, len(*listResp))
		result := SaveSignaturesOutput{
			WrongContent: &TruePtr,
			Details:      &details,
		}

		svc.logger.Debug(fmt.Sprintf("%s GetUpload() return", ctxId),
			zap.Reflect("result", result),
		)
		return &result, nil

	}

	sha256SumsKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS", *params.Namespace, *params.Type, *params.Version)
	sha256SumsData := []byte(*params.Sha256Sums)
	storageErr := svc.storageSvc.SaveObject(ctxId, sha256SumsKey, sha256SumsData)
	if storageErr != nil {
		return nil, svc.handleError(ctxId, storageErr)
	}
	svc.logger.Debug(fmt.Sprintf("%s provideres.SaveSignatures().SaveObject() sha256sums", ctxId),
		zap.String("key", sha256SumsKey),
	)

	sha256SumsSigKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig", *params.Namespace, *params.Type, *params.Version)
	sha256SumsSigData, base64Err := base64.StdEncoding.DecodeString(*params.Sha256SumsSig)
	if base64Err != nil {
		return nil, svc.handleError(ctxId, base64Err)
	}
	storageErr = svc.storageSvc.SaveObject(ctxId, sha256SumsSigKey, sha256SumsSigData)
	if storageErr != nil {
		return nil, svc.handleError(ctxId, storageErr)
	}
	svc.logger.Debug(fmt.Sprintf("%s provideres.SaveSignatures().SaveObject() sha256sums.sig", ctxId),
		zap.String("key", sha256SumsSigKey),
	)

	sha256SumsSigPubKey := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig.pub", *params.Namespace, *params.Type, *params.Version)
	sha256SumsSigPubData := []byte(*params.Sha256SumsSigPub)
	storageErr = svc.storageSvc.SaveObject(ctxId, sha256SumsSigPubKey, sha256SumsSigPubData)
	if storageErr != nil {
		return nil, svc.handleError(ctxId, storageErr)
	}
	svc.logger.Debug(fmt.Sprintf("%s provideres.SaveSignatures().SaveObject() sha256sums.sig.pub", ctxId),
		zap.String("key", sha256SumsSigPubKey),
	)

	sha256SumsSigKeyId := fmt.Sprintf("providers/%[1]s/%[2]s/%[3]s/terraform-provider-%[2]s_%[3]s_SHA256SUMS.sig.keyid", *params.Namespace, *params.Type, *params.Version)
	sha256SumsSigKeyIdData := []byte(*params.KeyId)
	storageErr = svc.storageSvc.SaveObject(ctxId, sha256SumsSigKeyId, sha256SumsSigKeyIdData)
	if storageErr != nil {
		return nil, svc.handleError(ctxId, storageErr)
	}
	svc.logger.Debug(fmt.Sprintf("%s provideres.SaveSignatures().SaveObject() sha256sums.sig.keyid", ctxId),
		zap.String("key", sha256SumsSigKeyId),
	)

	result := SaveSignaturesOutput{}

	svc.logger.Debug(fmt.Sprintf("%s GetUpload() return", ctxId),
		zap.Reflect("result", result),
	)
	return &result, nil
}

func (svc *Providers) handleError(ctxId string, err error) error {
	return errors.Wrapf(err, "%s providers error", ctxId)
}
