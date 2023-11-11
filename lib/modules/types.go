package modules

type (
	GetDownloadUrlOutput struct {
		Url          *string
		ModuleExists bool
	}

	ListModuleVersionsResponse struct {
		Modules      []ModuleVersions `json:"modules"`
		ModuleExists bool
	}

	ModuleVersions struct {
		Versions []Version `json:"versions"`
	}

	Version struct {
		Version string `json:"version"`
	}

	UploadModuleResponse struct {
		Url          *string `json:"url"`
		ModuleExists bool
	}

	ModuleParams struct {
		Namespace *string
		Name      *string
		Provider  *string
		Version   *string
	}
)
