package storage

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

type StorageProvider interface {
	Cp (CpParams) (string, error)
}
