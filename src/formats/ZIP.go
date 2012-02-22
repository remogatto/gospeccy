package formats

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
)

type ZipArchive struct {
	reader *zip.Reader
}

func (a *ZipArchive) Filenames() []string {
	filenames := make([]string, len(a.reader.File))
	for i, file := range a.reader.File {
		filenames[i] = file.Name
	}
	return filenames
}

func (a *ZipArchive) Read(fileIndex int) ([]byte, error) {
	if (fileIndex < 0) || (fileIndex >= len(a.reader.File)) {
		return nil, errors.New("invalid file index")
	}

	readCloser, err := a.reader.File[fileIndex].Open()
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	return ioutil.ReadAll(readCloser)
}

func ReadZip(r io.ReaderAt, size int64) (*ZipArchive, error) {
	reader, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}

	return &ZipArchive{reader}, nil
}

func ReadZipFile(filePath string) (*ZipArchive, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}

	return &ZipArchive{&reader.Reader}, nil
}
