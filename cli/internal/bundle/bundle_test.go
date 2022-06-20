package bundle

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestSyntheticTarBundling(t *testing.T) {
	dir := t.TempDir()
	// dir := "."
	fn := filepath.Join(dir, "output.tar")
	out, err := os.Create(fn)
	if err != nil {
		t.Fatal(err)
	}

	var files = []struct {
		Name string
	}{
		{"bundle_test.go"},
		{"bundle.go"},
	}

	t.Run("packing", func(t *testing.T) {
		tw := tar.NewWriter(out)

		for _, file := range files {

			f, err := os.Open(file.Name)
			if err != nil {
				t.Error(err)
			}
			info, err := f.Stat()

			if err != nil {
				t.Error(err)
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				t.Error(err)
			}

			header.Name = file.Name

			if err := tw.WriteHeader(header); err != nil {
				t.Error(err)
			}
			_, err = io.Copy(tw, f)
			if err != nil {
				t.Error(err)
			}
			err = f.Close()
			if err != nil {
				t.Error(err)
			}
		}
		if err := tw.Close(); err != nil {
			t.Error(err)
		}
	})

	if err := out.Close(); err != nil {
		t.Error(err)
	}

	rout, err := os.Open(fn)
	if err != nil {
		t.Error(err)
	}

	t.Run("unpacking", func(t *testing.T) {
		tr := tar.NewReader(rout)

		for {
			header, err := tr.Next()

			if err == io.EOF {
				break
			}
			fl := filepath.Join(dir, header.Name)

			f, err := os.Create(fl)
			if err != nil {
				t.Error(err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				t.Error(err)
			}
			if err = f.Close(); err != nil {
				t.Fatal(err)
			}
		}

		for _, f := range files {
			_, err := os.Stat(filepath.Join(dir, f.Name))
			if err != nil {
				t.Error(err)
			}
		}
	})
	if err = rout.Close(); err != nil {
		t.Error(err)
	}
}

func TestSynthenticTarGzipBundling(t *testing.T) {
	out, err := os.Create(filepath.Join(t.TempDir(), "output.tar.gz"))
	if err != nil {
		t.Fatal(err)
	}

	gw := gzip.NewWriter(out)
	tw := tar.NewWriter(gw)

	var files = []struct {
		Name string
	}{
		{"bundle_test.go"},
		{"bundle.go"},
	}
	for _, file := range files {

		f, err := os.Open(file.Name)
		if err != nil {
			t.Fatal(err)
		}
		info, err := f.Stat()

		if err != nil {
			t.Fatal(err)
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			t.Fatal(err)
		}

		header.Name = file.Name

		if err := tw.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		_, err = io.Copy(tw, f)
		if err != nil {
			t.Fatal(err)
		}
		err = f.Close()
		if err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := out.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestSynthenticZipBundling(t *testing.T) {
	// dir := t.TempDir()
	dir := "."
	out, err := os.Create(filepath.Join(dir, "output.zip"))
	if err != nil {
		t.Fatal(err)
	}

	zw := zip.NewWriter(out)

	files := []struct {
		Name, Body string
	}{
		{"demo/readme.txt", "This archive contains some text files."},
		{"demo/gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"demo/todo.txt", "Get animal handling licence.\nWrite more examples."},
	}

	for _, file := range files {
		f, err := zw.Create(file.Name)
		if err != nil {
			t.Error(err)
		}

		_, err = f.Write([]byte(file.Body))
		if err != nil {
			t.Error(err)
		}
	}

	err = zw.Close()
	if err != nil {
		t.Error(err)
	}
}
