package session

import (
	"fmt"
	"github.com/nbio/st"
	"testing"
)

func TestStripFormatting(t *testing.T) {
	text := "\x031,9\x02\x02<<SearchBot>> Your search for \"\x0312,9\x02\x02alias grace\x031,9\x02\x02\" returned 8 matches. Sending results to you as\x0312\x02\x02 SearchOok_results_for_ alias grace.txt.zip\x031,9\x02\x02. Search took 0.67 seconds."
	expected := "<<SearchBot>> Your search for \"alias grace\" returned 8 matches. Sending results to you as SearchOok_results_for_ alias grace.txt.zip. Search took 0.67 seconds."

	st.Expect(t, stripFormatting(text), expected)
}

func TestProcessSearchLine(t *testing.T) {
	// setup session for testing
	sess := NewSession("fuubar")
	sess.conn = new(ConnMock)
	sess.conn.on("Privmsg"
	query := "Japan at War"
	// TODO: stub sess.conn
	sess.searchBook(query)

	line := fmt.Sprintf("<<SearchBot>> Your search for \"%s\" has been accepted. Searching...", query)
	sess.processSearchLine(line)
}
