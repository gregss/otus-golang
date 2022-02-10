package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "/dev/urandom" {
		return ErrUnsupportedFile
	}

	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	defer fileFrom.Close()

	fstat, _ := fileFrom.Stat()
	if fstat.Mode().IsDir() {
		return errors.New("is directory")
	}

	fsize := fstat.Size()
	if fsize < offset {
		return ErrOffsetExceedsFileSize
	}

	if offset > 0 {
		fileFrom.Seek(offset, io.SeekStart)
	}

	fileTo, _ := os.Create(toPath)
	defer fileTo.Close()
	bufsize := 3
	buf := make([]byte, bufsize)
	bar := pb.StartNew(bufsize)

	sumreaded := 0
	for offset < fsize {
		read, err := fileFrom.Read(buf)
		sumreaded += read
		if limit > 0 && sumreaded > int(limit) {
			read -= (sumreaded - int(limit))
		}

		if err != nil {
			log.Panicf("failed to read: %v", err)
		}

		_, err = fileTo.Write(buf[:read])
		if err != nil {
			log.Panicf("failed to write: %v", err)
		}

		bar.Increment()

		if read < bufsize {
			break
		}

		offset += int64(read)
	}

	bar.Finish()

	return nil
}
