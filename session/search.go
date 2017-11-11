package session

import (
	"github.com/looplab/fsm"
)

// Search states
const (
	UNSTARTED                  = "UNSTARTED"
	SEARCH_SUBMITTED           = "SEARCH_SUBMITTED"
	SEARCH_ACCEPTED            = "SEARCH_ACCEPTED"
	SEARCH_RESULTS_FOUND       = "SEARCH_RESULTS_FOUND"
	SEARCH_RESULTS_NOT_FOUND   = "SEARCH_RESULTS_NOT_FOUND"
	DOWNLOAD_REQUEST_SUBMITTED = "DOWNLOAD_REQUEST_SUBMITTED"
	DOWNLOAD_REQUEST_ACCEPTED  = "DOWNLOAD_REQUEST_ACCEPTED"
	DOWNLOAD_SENT              = "DOWNLOAD_SENT"
)

// Search represents a book search, including the query and what
// stage the search is in
type Search struct {
	// query is the original search string e.g. 'Japan at War: An Oral History by Haruko Taya Cook'
	query string
	FSM   *fsm.FSM
}

// SearchRequest holds all the data needed to run a book search;
// used with workers to defer search request execution
type SearchRequest struct {
	query string
}

func NewSearch(query string) *Search {
	s := &Search{
		query: query,
	}

	s.FSM = fsm.NewFSM(
		UNSTARTED,
		fsm.Events{
			{Name: "submit", Src: []string{UNSTARTED}, Dst: SEARCH_SUBMITTED},
			{Name: "accept", Src: []string{SEARCH_SUBMITTED}, Dst: SEARCH_ACCEPTED},
			{Name: "foundResults", Src: []string{SEARCH_ACCEPTED}, Dst: SEARCH_RESULTS_FOUND},
			{Name: "noResults", Src: []string{SEARCH_ACCEPTED}, Dst: SEARCH_RESULTS_NOT_FOUND},
			{Name: "requestDownload", Src: []string{SEARCH_RESULTS_FOUND}, Dst: DOWNLOAD_REQUEST_SUBMITTED},
			{Name: "acceptDownload", Src: []string{DOWNLOAD_REQUEST_SUBMITTED}, Dst: DOWNLOAD_REQUEST_ACCEPTED},
			{Name: "sentDownload", Src: []string{DOWNLOAD_REQUEST_ACCEPTED}, Dst: DOWNLOAD_SENT},
		},
		fsm.Callbacks{},
	)

	return s
}

func (search *Search) submit() (err error) {
	return search.FSM.Event("submit")
}

func (search *Search) accept() (err error) {
	return search.FSM.Event("accept")
}

func (search *Search) foundResults() (err error) {
	return search.FSM.Event("foundResults")
}

func (search *Search) noResults() (err error) {
	return search.FSM.Event("noResults")
}

func (search *Search) requestDownload() (err error) {
	return search.FSM.Event("requestDownload")
}

func (search *Search) acceptDownload() (err error) {
	return search.FSM.Event("acceptDownload")
}

func (search *Search) sendDownload() (err error) {
	return search.FSM.Event("sendDownload")
}
