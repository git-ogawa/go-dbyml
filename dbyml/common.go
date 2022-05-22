package dbyml

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

func GetTarContext(dockerFile string) (buf *bytes.Buffer) {
	buf = new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	abs, _ := filepath.Abs(dockerFile)
	dockerFileReader, err := os.Open(abs)
	if err != nil {
		log.Fatal(err, " :unable to open Dockerfile")
	}
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		log.Fatal(err, " :unable to read dockerfile")
	}

	tarHeader := &tar.Header{
		Name: dockerFile,
		Size: int64(len(readDockerFile)),
	}
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		log.Fatal(err, " :unable to write tar header")
	}
	_, err = tw.Write(readDockerFile)
	if err != nil {
		log.Fatal(err, " :unable to write tar body")
	}
	return buf
}

// Centering returns the string centered within the specified length.
func Centering(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}

// PrintCenter shows centered string with padding specified character on each side of the string to stdout.
func PrintCenter(s string, w int, padding string) {
	side := strings.Repeat(padding, w)
	center := Centering(s, w)
	fmt.Printf("%v%v%v\n", side, center, side)
}

func ShowArrayElement(name string, values []string) {
	for i, v := range values {
		if i == 0 {
			fmt.Printf("%-30v: %v\n", name, v)
		} else {
			fmt.Printf("%-30v: %v\n", "", v)
		}
	}
}

func ShowMapElement(name string, iter *reflect.MapIter) {
	cnt := 0
	for iter.Next() {
		if iter.Value().Kind() == reflect.Ptr {
			if cnt == 0 {
				fmt.Printf("%-30v: %v: %v\n", name, iter.Key(), iter.Value().Elem())
			} else {
				fmt.Printf("%-30v: %v: %v\n", "", iter.Key(), iter.Value().Elem())
			}
		} else {
			if cnt == 0 {
				fmt.Printf("%-30v: %v: %v\n", name, iter.Key(), iter.Value())
			} else {
				fmt.Printf("%-30v: %v: %v\n", "", iter.Key(), iter.Value())
			}
		}
		cnt += 1
	}
}

func splitString(s string) (k, v string) {
	arr := strings.Split(s, ":")
	if len(arr) == 2 {
		return strings.TrimSpace(arr[0]), strings.TrimSpace(arr[1])
	} else {
		panic("this is no")
	}
}

func splitStringPointer(s string) (k string, v *string) {
	arr := strings.Split(s, ":")
	if len(arr) == 2 {
		ptr := strings.TrimSpace(arr[1])
		return strings.TrimSpace(arr[0]), &ptr
	} else {
		panic("this is no")
	}
}