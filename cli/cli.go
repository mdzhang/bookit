package cli

import (
	"github.com/mdzhang/bookit/logger"
	"github.com/mdzhang/bookit/session"
	"github.com/mdzhang/bookit/version"
	"gopkg.in/alecthomas/kingpin.v2"
	"os"
)

// CLI flags
var (
	nick  = kingpin.Flag("nick", "IRC nick to use").Short('n').Required().String()
	query = kingpin.Flag("query", "Book to query for e.g. 'Japan at War: An Oral History by Haruko Taya Cook'").Short('q').Required().String()
)

// Run CLI parser and bookit program
func Run() {
	logger.Init()
	kingpin.Version("v" + version.Version)

	kingpin.Parse()

	logger.Info("Starting session...")
	session := session.NewSession(*nick)
	quit := make(chan bool)

	if err := session.Connect(quit); err != nil {
		logger.Info("Connection error: %s\n", err.Error())
		os.Exit(1)
	}

	logger.Info("Searching book...")
	session.SearchBook(*query)
	<-quit
}
