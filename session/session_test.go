package session

import (
	"github.com/nbio/st"
	"testing"
)

func TestStripFormatting(t *testing.T) {
	text := "1,9<<SearchBot>> Your search for \"12,9alias grace1,9\" returned 8 matches. Sending results to you as12 SearchOok_results_for_ alias grace.txt.zip1,9. Search took 0.67 seconds."
	expected := "<<SearchBot>> Your search for \"alias grace\" returned 8 matches. Sending results to you as SearchOok_results_for_ alias grace.txt.zip. Search took 0.67 seconds."

	st.Expect(t, stripFormatting(text), expected)
}
