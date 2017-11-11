package cli

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/mdzhang/bookit/version"
	"gopkg.in/alecthomas/kingpin.v2"
)

// CLI flags
var (
	nick = kingpin.Flag("nick", "IRC nick to use").String()
)

// Run CLI parser and bookit proram
func Run() {
	kingpin.Version("v" + version.Version)

	kingpin.Parse()

	fmt.Printf("%s => Using nick: '%s'", color.GreenString("bookit"), *nick)
}
