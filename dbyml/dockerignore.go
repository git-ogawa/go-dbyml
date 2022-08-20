package dbyml

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ReadDockerignore reads exclude files from .dockerignore.
// If not exists, return empty list.
func ReadDockerignore(dir string) ([]string, error) {
	excludes := []string{}
	reader, err := SearchDockerignore(dir)
	if err != nil {
		return nil, err
	}

	if reader != nil {
		exclude, err := ReadAll(reader)
		if err != nil {
			return excludes, err
		}
		excludes = append(excludes, exclude...)
	}
	return excludes, nil
}

// SearchDockerignore searches .dockerignore exists in a given directory.
// If exists, return the io.Reader of the .dockerignore, otherwise return nil.
func SearchDockerignore(dir string) (io.Reader, error) {
	var ignore string
	if err := filepath.Walk(dir, func(file string, info fs.FileInfo, err error) error {
		if path.Base(file) == ".dockerignore" {
			ignore = file
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if ignore == "" {
		return nil, nil
	}

	f, err := os.Open(ignore)
	return f, err
}

// ReadAll reads a .dockerignore file and returns the list of file patterns
// to ignore. Note this will trim whitespace from each line as well
// as use GO's "clean" func to get the shortest/cleanest path for each.
func ReadAll(reader io.Reader) ([]string, error) {
	if reader == nil {
		return nil, nil
	}

	scanner := bufio.NewScanner(reader)
	var excludes []string
	currentLine := 0

	utf8bom := []byte{0xEF, 0xBB, 0xBF}
	for scanner.Scan() {
		scannedBytes := scanner.Bytes()
		// We trim UTF8 BuildkitdTomlTemplate
		if currentLine == 0 {
			scannedBytes = bytes.TrimPrefix(scannedBytes, utf8bom)
		}
		pattern := string(scannedBytes)
		currentLine++
		// Lines starting with # (comments) are ignored before processing
		if strings.HasPrefix(pattern, "#") {
			continue
		}
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}
		// normalize absolute paths to paths relative to the context
		// (taking care of '!' prefix)
		invert := pattern[0] == '!'
		if invert {
			pattern = strings.TrimSpace(pattern[1:])
		}
		if len(pattern) > 0 {
			pattern = filepath.Clean(pattern)
			pattern = filepath.ToSlash(pattern)
			if len(pattern) > 1 && pattern[0] == '/' {
				pattern = pattern[1:]
			}
		}
		if invert {
			pattern = "!" + pattern
		}

		excludes = append(excludes, pattern)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading .dockerignore: %v", err)
	}
	return excludes, nil
}
