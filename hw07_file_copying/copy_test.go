package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("error. not supported file", func(t *testing.T) {
		err := Copy("/dev/urandom", "/", 0, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "unsupported file")
	})

	t.Run("error. is directory", func(t *testing.T) {
		err := Copy("testdata", "./", 0, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "is directory")
	})

	t.Run("error. offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/out_offset0_limit0.txt", "./", 100000, 0)
		require.Error(t, err)
		require.Equal(t, err.Error(), "offset exceeds file size")
	})

	t.Run("success", func(t *testing.T) {
		fname1 := "testfile"
		fname2 := "testfile2"
		file, _ := os.Create(fname1)
		defer os.Remove(fname1)
		defer os.Remove(fname2)

		file.Write([]byte("0123456789"))
		require.Nil(t, Copy(fname1, fname2, 0, 0))

		fileFrom1, _ := os.Open(fname1)
		fileFrom2, _ := os.Open(fname2)
		defer fileFrom1.Close()
		defer fileFrom2.Close()
		fstat1, _ := fileFrom1.Stat()
		fstat2, _ := fileFrom2.Stat()
		require.Equal(t, fstat1.Size(), fstat2.Size())
	})

	t.Run("success", func(t *testing.T) {
		fname1 := "testfile"
		fname2 := "testfile2"
		file, _ := os.Create(fname1)
		defer os.Remove(fname1)
		defer os.Remove(fname2)
		file.Write([]byte("0123456789"))
		require.Nil(t, Copy(fname1, fname2, 0, 0))

		fileFrom2, _ := os.Open(fname2)
		defer fileFrom2.Close()

		fstat2, _ := fileFrom2.Stat()
		require.Equal(t, int64(10), fstat2.Size())

		Copy(fname1, fname2, 5, 0)
		fstat2, _ = fileFrom2.Stat()
		require.Equal(t, int64(5), fstat2.Size())

		Copy(fname1, fname2, 5, 3)
		fstat2, _ = fileFrom2.Stat()
		require.Equal(t, int64(3), fstat2.Size())

		Copy(fname1, fname2, 0, 20)
		fstat2, _ = fileFrom2.Stat()
		require.Equal(t, int64(10), fstat2.Size())
	})
}
