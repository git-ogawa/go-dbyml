package dbyml

import (
	"fmt"
	"reflect"
	"strings"
)

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

func showMapElement(name string, iter *reflect.MapIter) {
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
		cnt++
	}
}
