package update

import (
	"compress/bzip2"
	"io"
	"os"

	"github.com/jaczerob/clerk-gopher/internal/util"
)

func NewUpdateFile(name, path, url string) (f *UpdateFile, err error) {
	exists := util.FileExists(path)

	var hash string
	var out *os.File
	if exists {
		hash, err = util.GetHash(path)
		if err != nil {
			return
		}

		out, err = util.GetFile(path)
		if err != nil {
			return
		}
	} else {
		hash = ""
		out = nil
	}

	f = &UpdateFile{
		Name:   name,
		Path:   path,
		URL:    url,
		Hash:   hash,
		Exists: exists,
		out:    out,
	}

	return
}

func (f *UpdateFile) Read(body io.ReadCloser) (err error) {
	defer f.out.Close()
	decompressed := bzip2.NewReader(body)
	_, err = io.Copy(f.out, decompressed)
	return
}
