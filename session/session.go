package session

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/mdzhang/bookit/logger"
	"os"
	"regexp"
)

// Session represents an IRC connection, along
// with a queue for deferring client interactions
type Session struct {
	conn      *irc.Conn
	connected chan bool
	queue     chan SearchRequest
	search    *Search
	nick      string
	server    string
	channel   string
}

// NewSession creates a Session with reasonable defaults
func NewSession(nick string) *Session {
	sess := &Session{
		connected: make(chan bool, 1),
		queue:     make(chan SearchRequest),
		nick:      nick,
		server:    "irc.irchighway.net:6667",
		channel:   "#ebooks",
	}
	return sess
}

func (sess *Session) Connect(quit chan bool) (err error) {
	cfg := irc.NewConfig(sess.nick)
	cfg.Server = sess.server
	c := irc.Client(cfg)

	c.HandleFunc(irc.CONNECTED, sess.h_joinChannel)

	for _, cmd := range [...]string{irc.PRIVMSG, irc.NOTICE} {
		c.HandleFunc(cmd, sess.h_processLine)
	}

	c.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { quit <- true })

	sess.conn = c

	go sess.startWorker()

	logger.Info("Connecting to %s", sess.server)
	return c.Connect()
}

func (sess *Session) SearchBook(query string) {
	logger.Info("Queued search request '%s'", query)
	request := SearchRequest{query: query}
	sess.queue <- request
}

func (sess *Session) searchBook(query string) {
	search := NewSearch(query)
	sess.search = search
	s := fmt.Sprintf("@search %s", query)
	logger.Info("Sending message %s", s)
	sess.conn.Privmsg(sess.channel, s)
	search.submit()
}

func (sess *Session) startWorker() {
	logger.Info("[worker] Waiting until connected...")
	<-sess.connected

	logger.Info("[worker] Connected. Processing requests...")
	for {
		request := <-sess.queue
		sess.searchBook(request.query)
	}
}

func (sess *Session) h_joinChannel(conn *irc.Conn, line *irc.Line) {
	logger.Info("Joining %s", sess.channel)
	conn.Join(sess.channel)
	sess.connected <- true
}

func (sess *Session) h_processLine(conn *irc.Conn, line *irc.Line) {
	logLine := func(line *irc.Line) {
		logger.Info("[irc] [%s] %s", line.Nick, line.Text())
	}

	if line.Target() != sess.nick {
		logger.Info("[irc] --")
		return
	}

	switch nick := line.Nick; nick {
	case "SearchOok":
		logLine(line)
		sess.processSearchLine(line.Text())
	case "ChanServ", sess.nick:
		logLine(line)
	default:
		logger.Info("[irc] --")
	}
}

// stripFormatting removes IRC text color formatting
//	see https://en.wikichip.org/wiki/irc/colors
func stripFormatting(text string) string {
	r := regexp.MustCompile("\x03\\d{0,2}(,\\d{0,2})?(\x02\x02)?")
	return r.ReplaceAllLiteralString(text, "")
}

// TODO: refactor candidate
func (sess *Session) processSearchLine(line string) {
	text := stripFormatting(line)

	searchAccepted := fmt.Sprintf("<<SearchBot>> Your search for \"%s\" has been accepted. Searching...", sess.search.query)
	searchAcceptedRegex := regexp.MustCompile(searchAccepted)

	foundResults := fmt.Sprintf("<<SearchBot>> Your search for \"%s\" returned 8 matches. Sending results to you as SearchOok_results_for_ %s.txt.zip. Search took \\d+\\.\\d+ seconds.", sess.search.query, sess.search.query)
	foundResultsRegex := regexp.MustCompile(foundResults)

	noResults := fmt.Sprintf("Sorry, your search for \"%s\" returned no matches.*", sess.search.query)
	noResultsRegex := regexp.MustCompile(noResults)

	resultsFileName := fmt.Sprintf("SearchOok_results_for_ %s.txt.zip", sess.search.query)
	dccSent := fmt.Sprintf("DCC Send %s .*", resultsFileName)
	dccSentRegex := regexp.MustCompile(dccSent)

	fileRequestAccepted := fmt.Sprintf(".*Request Accepted ? File: %s", sess.search.requestedFile)
	fileRequestAcceptedRegex := regexp.MustCompile(fileRequestAccepted)

	fileReceived := fmt.Sprintf("DCC Send %s .*", sess.search.requestedFile)
	fileReceivedRegex := regexp.MustCompile(fileReceived)

	if searchAcceptedRegex.MatchString(text) {
		logger.Info("Search accepted")
		sess.search.accept()
	} else if foundResultsRegex.MatchString(text) {
		logger.Info("Found results")
		sess.search.foundResults()
	} else if dccSentRegex.MatchString(text) {
		logger.Info("Search results received")
		sess.search.dccSent()
		// TODO: download the DCC file (resultsFileName) to disk
		//			 read and choose file source
		//			 enqueue request to download file and update search state
	} else if noResultsRegex.MatchString(text) {
		logger.Info("No results")
		sess.search.noResults()
		// TODO: prob best to propagate an error here and caller can decide
		//			 what to do with it
		os.Exit(1)
	} else if fileRequestAcceptedRegex.MatchString(text) {
		logger.Info("File request accepted")
		sess.search.fileRequestAccepted()
	} else if fileReceivedRegex.MatchString(text) {
		logger.Info("File received")
		// TODO: download the DCC file to disk
		// TODO: prob best to report done some other way
		os.Exit(1)
	}

	logger.Info("Processed text: %s", text)
}
