package whizz_cli

import (
	"fmt"
	"reflect"

	"github.com/infra-whizz/wzlib"
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

func (wcf *WhizzCliFormatter) ListSystems(systems []interface{}) {
	statuses := map[int]string{
		wzlib.CLIENT_STATUS_ACCEPTED: "Accepted",
		wzlib.CLIENT_STATUS_NEW:      "New",
		wzlib.CLIENT_STATUS_REJECTED: "Rejected",
	}
	if len(systems) > 0 {
		fmt.Printf("%d machine(s) found:\n", len(systems))
		for idx, system := range systems {
			s := system.(map[string]interface{})
			fmt.Printf("  %d.  %s\t%s\t%s\n", idx+1, s["RsaFp"].(string)[:8], statuses[int(s["Status"].(int64))], s["Fqdn"])
		}
	} else {
		fmt.Println("No machines found by this criteria")
	}

}
