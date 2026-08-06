package archive

import (
	"archive/tar"
	"bytes"
	"io"
	"os"
)

func Tar(dir string) (io.Reader, error) {
	var buff bytes.Buffer

	wrt := tar.NewWriter(&buff)
	if err := wrt.AddFS(os.DirFS(dir)); err != nil {
		return nil, err
	}

	return &buff, wrt.Close()
}
