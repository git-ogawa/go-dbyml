package dbyml

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/moby/moby/pkg/fileutils"
)

// GetBuildContext makes tar archive of files and directories in a given directory, and returns
// byte.Buffer of the archive. The buffer is used for build context to build an image.
func GetBuildContext(dir string) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	excludes, err := ReadDockerignore(dir)
	if err != nil {
		panic(err)
	}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Check exclude
		if rm, _ := IsExclude(path, excludes); rm {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		// Write header
		if err := tw.WriteHeader(&tar.Header{
			Name:    path,
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		}); err != nil {
			return err
		}

		// Write body
		_, err = tw.Write(b)
		if err != nil {
			log.Fatal(err, " :unable to write tar body")
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return buf
}

// IsExclude returns true if file matches any of the patterns and isn't excluded by any of the subsequent patterns.
func IsExclude(file string, exclude []string) (bool, error) {
	return fileutils.Matches(file, exclude)
}

// GetBuildContext makes tar archive of files and directories in a given directory, and returns
// byte.Buffer of the archive. The buffer is used for build context to build an image in bubildkitd container.
func GetBuildkitContext(dir string) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	excludes, err := ReadDockerignore(dir)
	if err != nil {
		panic(err)
	}

	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Check exclude
		if rm, _ := IsExclude(path, excludes); rm {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		b, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}

		// Replace directory
		str := regexp.MustCompile(dir)
		path = str.ReplaceAllString(path, ".")

		// Write header
		if err := tw.WriteHeader(&tar.Header{
			Name:    path,
			Mode:    int64(info.Mode()),
			ModTime: info.ModTime(),
			Size:    info.Size(),
		}); err != nil {
			return err
		}

		// Write body
		_, err = tw.Write(b)
		if err != nil {
			log.Fatal(err, " :unable to write tar body")
		}
		return nil
	}); err != nil {
		panic(err)
	}
	return buf
}
