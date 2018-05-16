// Package agrovoc offers a simple way to search for a concept on agrovoc.
package agrovoc

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	endpoint = "http://artemide.art.uniroma2.it:8081/agrovoc/rest/v1/"
)

func search(r SearchRequest) (*SearchResults, error) {
	var surl strings.Builder
	surl.WriteString(endpoint + "search")

	surl.WriteString("?query=" + url.QueryEscape(r.Query))
	surl.WriteString("&lang=" + r.Lang)
	surl.WriteString("&labellang=" + r.Labellang)
	surl.WriteString("&vocab=" + r.Vocab)
	surl.WriteString("&type=" + r.Type)
	surl.WriteString("&parent=" + r.Parent)
	surl.WriteString("&group=" + r.Group)
	surl.WriteString("&maxhits=" + strconv.Itoa(r.Maxhits))
	surl.WriteString("&offset=" + strconv.Itoa(r.Offset))
	surl.WriteString("&fields=" + strings.Join(r.Fields, " "))
	surl.WriteString("&unique=" + strconv.FormatBool(r.Unique))

	res, err := http.Get(surl.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	searchResults := new(SearchResults)
	err = json.NewDecoder(res.Body).Decode(searchResults)
	if err != nil {
		return nil, err
	}

	return searchResults, nil
}

// ExactMatch takes a query as parameter, performs a search on agrovoc and returns a uri/label couple list.
func ExactMatch(query string) ([]Result, error) {
	searchResults, err := search(SearchRequest{
		Query:     query,
		Labellang: "fr",
		Vocab:     "agrovoc",
	})
	if err != nil {
		return nil, err
	}

	var results []Result

	for _, result := range searchResults.Results {
		results = append(results, Result{result.PrefLabel, result.URI})
	}

	return results, nil
}
