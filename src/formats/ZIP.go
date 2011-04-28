package formats

import (
	"archive/zip"
//	"io"
	"io/ioutil"
	"os"
)

type ZipArchive struct {
	reader *zip.ReadCloser
}

func (a *ZipArchive) Filenames() []string {
	filenames := make([]string, len(a.reader.File))
	for i, file := range a.reader.File {
		filenames[i] = file.Name
	}
	return filenames
}

func (a *ZipArchive) Read(fileIndex int) ([]byte, os.Error) {
	if (fileIndex < 0) || (fileIndex >= len(a.reader.File)) {
		return nil, os.NewError("invalid file index")
	}

	readCloser, err := a.reader.File[fileIndex].Open()
	if err != nil {
		return nil, err
	}

	defer readCloser.Close()

	return ioutil.ReadAll(readCloser)
}

// func ReadZip(r io.ReaderAt, size int64) (*ZipArchive, os.Error) {
// 	reader, err := zip.NewReader(r, size)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &ZipArchive{reader}, nil
// }

func ReadZipFile(filePath string) (*ZipArchive, os.Error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}

	return &ZipArchive{reader}, nil
}
