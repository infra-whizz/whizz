package whizz_cli

import (
	"fmt"
	"reflect"
)

type WhizzCliFormatter struct{}

func NewWhizzCliFormatter() *WhizzCliFormatter {
	wcf := new(WhizzCliFormatter)
	return wcf
}

// Map just dump a content of a map of string/interface to the STDOUT, keys sorted
func (wcf *WhizzCliFormatter) Map(data map[string]interface{}) {
	wcf.printMap(data, "")
}

func (wcf *WhizzCliFormatter) printMap(data map[string]interface{}, ident string) {
	fmt.Println(ident + "|")
	for k, v := range data {
		if reflect.TypeOf(v).Kind() == reflect.Map {
			wcf.printMap(v.(map[string]interface{}), ident+"  ")
		} else {
			fmt.Printf("%s+-%s: %v\n", ident, k, v)
		}
	}
}

// HostnameWithFp prints hostname and fingerprint underneath
func (wcf *WhizzCliFormatter) HostnameWithFp(idx int, hostname string, fingerprint string) {
	fmt.Printf("%d.  %s\n    +- %s\n", idx, hostname, fingerprint)
}
