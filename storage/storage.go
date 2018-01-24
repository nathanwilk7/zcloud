package storage

type StorageProvider interface {
	Upload (UploadParams) (string, error)
	Download (DownloadParams) (string, error)
	Ls (LsParams) (string, error)
	Rm (RmParams) (string, error)
	Mv (MvParams) (string, error)
	Mb (MbParams) (string, error)
	Rb (RmParams) (string, error)
	Sync (SyncParams) (string, error)
	StorageURLPrefixReplacement() string
}

type UploadParams struct {
	Src, Dest string
	Recursive bool
}

func NewUploadParams (src string, dest string) UploadParams {
	return UploadParams{
		Src: src,
		Dest: dest,
		Recursive: false,
	}
}

type DownloadParams struct {
	Src, Dest string
	Recursive bool
}

func NewDownloadParams (src string, dest string) DownloadParams {
	return DownloadParams{
		Src: src,
		Dest: dest,
		Recursive: false,
	}
}

type LsParams struct {
	Url string
	Recursive bool
}

func NewLsParams (url string) LsParams {
	return LsParams{
		Url: url,
		Recursive: false,
	}
}

type RmParams struct {
	Url string
	Recursive bool
}

func NewRmParams (url string) RmParams {
	return RmParams{
		Url: url,
		Recursive: false,
	}
}

type MvParams struct {
	Url string
	Recursive bool
}

func NewMvParams (url string) MvParams {
	return MvParams{
		Url: url,
		Recursive: false,
	}
}

type MbParams struct {
	Url string
}

func NewMbParams (url string) MbParams {
	return MbParams{
		Url: url,
	}
}

type RbParams struct {
	Url string
	Force bool
}

func NewRbParams (url string) RbParams {
	return RbParams{
		Url: url,
		Force: false,
	}
}

type SyncParams struct {
	Url string
	Delete bool
}

func NewSyncParams (url string) SyncParams {
	return SyncParams{
		Url: url,
		Delete: false,
	}
}
