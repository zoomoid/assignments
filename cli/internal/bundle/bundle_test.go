package bundle

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/zoomoid/assignments/v1/internal/context"
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
	dir := t.TempDir()
	fn := filepath.Join(dir, "output.tar.gz")
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
		gw := gzip.NewWriter(out)
		tw := tar.NewWriter(gw)
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
	})

	if err := out.Close(); err != nil {
		t.Error(err)
	}

	rout, err := os.Open(fn)
	if err != nil {
		t.Error(err)
	}

	t.Run("unpacking", func(t *testing.T) {
		gr, err := gzip.NewReader(rout)
		if err != nil {
			t.Fatal(err)
		}
		tr := tar.NewReader(gr)

		for {
			header, err := tr.Next()

			if err == io.EOF {
				break
			}
			fl := filepath.Join(dir, header.Name)

			f, err := os.Create(fl)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(f, tr); err != nil {
				t.Fatal(err)
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

func TestSynthenticZipBundling(t *testing.T) {
	dir := t.TempDir()
	// dir := "."
	fn := filepath.Join(dir, "output.zip")
	out, err := os.Create(fn)
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

func TestMakeBundler(t *testing.T) {
	ctx, _ := context.NewDevelopment()
	bundler, err := New(ctx, &BundlerOptions{
		Target: "assignment-01.pdf",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Run("backend=BundlerBackendTar", func(t *testing.T) {
		bundler.Backend = BundlerBackendTar
		b, err := bundler.makeBundler()
		if err != nil {
			t.Fatal(err)
		}
		if b.Type() != BundlerBackendTar {
			t.Errorf("expected type of bundler to be %s, found %s", BundlerBackendTar, b.Type())
		}
	})
	t.Run("backend=BundlerBackendZip", func(t *testing.T) {
		bundler.Backend = BundlerBackendZip
		b, err := bundler.makeBundler()
		if err != nil {
			t.Fatal(err)
		}
		if b.Type() != BundlerBackendZip {
			t.Errorf("expected type of bundler to be %s, found %s", BundlerBackendZip, b.Type())
		}
	})
	t.Run("backend=BundlerBackendTarGzip", func(t *testing.T) {
		bundler.Backend = BundlerBackendTarGzip
		b, err := bundler.makeBundler()
		if err != nil {
			t.Fatal(err)
		}
		if b.Type() != BundlerBackendTarGzip {
			t.Errorf("expected type of bundler to be %s, found %s", BundlerBackendTarGzip, b.Type())
		}
	})
}
