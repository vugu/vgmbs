// Package vgmbs embeds SASS files provides a way to use them in your project.
package vgmbs

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// NOTE: uses https://github.com/shurcooL/vfsgen for embedding
// go get -u github.com/shurcooL/vfsgen/cmd/vfsgendev

//go:generate go run build.go

// NewFileWriter returns a FileWriter.
func NewFileWriter(dir string) *FileWriter {
	fw := &FileWriter{dir: dir}
	return fw
}

// FileWriter can write the SCSS files out to a directory (source comes from embedded files in this package).
type FileWriter struct {
	dir       string
	noVerify  bool
	overwrite bool
}

//// NoVerify disable file verification and just ignore if file already exists.
//func (fw *FileWriter) NoVerify() *FileWriter {
//	fw.noVerify = true
//	return fw
//}
//
//// Overwrite will cause any existing file to be overwritten.
//func (fw *FileWriter) Overwrite() *FileWriter {
//	fw.overwrite = true
//	return fw
//}

// TODO: NoMkdir option, MkdirAll option - or think about the cases here - probably
// "dist/out" should be willing to create both of those folders, but "/project/dist/out"
// should probably refuse to create anything except "out", because we don't really know
// the significance of an absolute path.

// Write will write the embedded SASS files in this packge to the specified directory.
func (fw *FileWriter) Write() error {
	return fw.writeAssetDir("/")
}

func (fw *FileWriter) writeAssetDir(assetDir string) error {

	err := os.Mkdir(filepath.Join(fw.dir, assetDir), 0755)
	if os.IsExist(err) {
		err = nil
	}
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	f1, err := assets.Open(assetDir)
	if err != nil {
		return fmt.Errorf("failed to open internal asset dir: %w", err)
	}
	defer f1.Close()
	fis, err := f1.Readdir(-1)
	if err != nil {
		return fmt.Errorf("failed to Readdir on internal asset dir: %w", err)
	}
	for _, fi := range fis {
		fi := fi
		err := func() error {

			if fi.IsDir() {
				err := fw.writeAssetDir(path.Join(assetDir, fi.Name()))
				if err != nil {
					return err
				}
				return nil
			}

			f, err := assets.Open(path.Join(assetDir, fi.Name()))
			if err != nil {
				return fmt.Errorf("failed to Open %q internal asset dir: %w", fi.Name(), err)
			}
			defer f.Close()
			outPath := filepath.Join(fw.dir, assetDir, fi.Name())
			fout, err := os.Create(outPath)
			if err != nil {
				return fmt.Errorf("failed to Create %q: %w", outPath, err)
			}
			defer fout.Close()
			_, err = io.Copy(fout, f)
			if err != nil {
				return fmt.Errorf("failed to io.Copy for %q: %w", outPath, err)
			}
			return nil
		}()
		if err != nil {
			return err
		}

	}

	return nil
}

// MustWrite is like Write but panics upon error.
func (fw *FileWriter) MustWrite() {
	err := fw.Write()
	if err != nil {
		panic(err)
	}
}

// TODO: method to return the list of files names in the correct sequence
// TODO: consider prefixing the files with "material_NN_", so it's clear when sitting with other files and sequence is indicated.
// e.g. "material_01_reset.scss" - it's pretty obvious what that is, even when sitting next to other files, and then
// "material_02_globals.scss" - it's clear what the sequence is on that.
