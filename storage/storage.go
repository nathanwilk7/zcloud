package storage

type StorageProvider interface {
	Cp (CpParams) (string, error)
	Ls (LsParams) (string, error)
	Rm (RmParams) (string, error)
	Mv (MvParams) (string, error)
	Mb (MbParams) (string, error)
	Rb (RmParams) (string, error)
	Sync (SyncParams) (string, error)
}

type CpParams struct {
	Src, Dest string
	Recursive bool
}

func NewCpParams (src string, dest string) CpParams {
	return CpParams{
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
