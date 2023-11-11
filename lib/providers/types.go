package providers

type (
	SaveSignaturesInput struct {
		Namespace        *string
		Type             *string
		Version          *string
		KeyId            *string
		Sha256Sums       *string
		Sha256SumsSig    *string
		Sha256SumsSigPub *string
	}

	SaveSignaturesOutput struct {
		ProviderExists *bool
		MetadataExists *bool
		WrongContent   *bool
		Details        *string
	}
)
