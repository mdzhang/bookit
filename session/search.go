package session

import (
	"github.com/looplab/fsm"
)

// Search states
const (
	// haven't started search
	UNSTARTED = "UNSTARTED"
	// asked for names of files matching search
	SEARCH_SUBMITTED = "SEARCH_SUBMITTED"
	// server accepted search
	SEARCH_ACCEPTED = "SEARCH_ACCEPTED"
	// server found names of file matching search
	SEARCH_RESULTS_FOUND = "SEARCH_RESULTS_FOUND"
	// server send file containing names of files matching search
	DCC_SENT = "DCC_SENT"
	// server found no files matching search
	SEARCH_RESULTS_NOT_FOUND = "SEARCH_RESULTS_NOT_FOUND"
	// client asked for specific file
	DOWNLOAD_REQUEST_SUBMITTED = "DOWNLOAD_REQUEST_SUBMITTED"
	// server received client request for specific file
	DOWNLOAD_REQUEST_ACCEPTED = "DOWNLOAD_REQUEST_ACCEPTED"
	// server sent file requested
	DOWNLOAD_SENT = "DOWNLOAD_SENT"
)

// Search represents a book search, including the query and what
// stage the search is in
type Search struct {
	// query is the original search string e.g. 'Japan at War: An Oral History by Haruko Taya Cook'
	query string
	// file requested; will be one of the results returned from server for query
	requestedFile string
	FSM           *fsm.FSM
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
			{Name: "dccSent", Src: []string{SEARCH_RESULTS_FOUND}, Dst: DCC_SENT},
			{Name: "noResults", Src: []string{SEARCH_ACCEPTED}, Dst: SEARCH_RESULTS_NOT_FOUND},
			{Name: "requestFile", Src: []string{DCC_SENT}, Dst: DOWNLOAD_REQUEST_SUBMITTED},
			{Name: "fileRequestAccepted", Src: []string{DOWNLOAD_REQUEST_SUBMITTED}, Dst: DOWNLOAD_REQUEST_ACCEPTED},
			{Name: "fileReceived", Src: []string{DOWNLOAD_REQUEST_ACCEPTED}, Dst: DOWNLOAD_SENT},
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

func (search *Search) dccSent() (err error) {
	return search.FSM.Event("dccSent")
}

func (search *Search) noResults() (err error) {
	return search.FSM.Event("noResults")
}

func (search *Search) requestFile() (err error) {
	return search.FSM.Event("requestFile")
}

func (search *Search) fileRequestAccepted() (err error) {
	return search.FSM.Event("fileRequestAccepted")
}

func (search *Search) sendDownload() (err error) {
	return search.FSM.Event("sendDownload")
}
