package session

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/mdzhang/bookit/logger"
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
		sess.processSearchLine(line)
	case "ChanServ", sess.nick:
		logLine(line)
	default:
		logger.Info("[irc] --")
	}
}

// stripFormatting removes IRC text formatting that looks like e.g. '12,9'
//	see https://en.wikichip.org/wiki/irc/colors
func stripFormatting(text string) string {
	r := regexp.MustCompile("\\d{1,2}(,\\d{1,2})?")

	return r.ReplaceAllLiteralString(text, "")
}

func (sess *Session) processSearchLine(line *irc.Line) {
	text := stripFormatting(line.Text())

	// <<SearchBot>> Your search for "alias grace" has been accepted. Searching...

	// <<SearchBot>> Your search for "alias grace" returned 8 matches. Sending results to you as SearchOok_results_for_ alias grace.txt.zip. Search took 0.67 seconds.
	// DCC Send SearchOok_results_for_ alias grace.txt.zip CRC(838E30AB) (unseen.edu)

	// 01:48 pondering42 Notice:  Request Accepted ? File: Margaret Atwood - Alias Grace (v5.0) (epub).rar ? Queue Position: 187 ? Allowed: 1 of 32 ? Min CPS: 50 ? OmenServe v2.72 ?

	logger.Info("Processed text: %s", text)
}
