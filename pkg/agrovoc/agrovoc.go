// Package agrovoc offers a simple way to search for a concept on agrovoc.
package agrovoc

import (
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
)

var (
	endpoint = "http://artemide.art.uniroma2.it:8081/agrovoc/rest/v1/"
)

func search(r SearchRequest) (*SearchResults, error) {
	var url strings.Builder
	url.WriteString(endpoint + "search")

	url.WriteString("?query=" + r.Query)
	url.WriteString("&lang=" + r.Lang)
	url.WriteString("&labellang=" + r.Labellang)
	url.WriteString("&vocab=" + r.Vocab)
	url.WriteString("&type=" + r.Type)
	url.WriteString("&parent=" + r.Parent)
	url.WriteString("&group=" + r.Group)
	url.WriteString("&maxhits=" + strconv.Itoa(r.Maxhits))
	url.WriteString("&offset=" + strconv.Itoa(r.Offset))
	url.WriteString("&fields=" + strings.Join(r.Fields, " "))
	url.WriteString("&unique=" + strconv.FormatBool(r.Unique))

	res, err := http.Get(url.String())
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

	results := make([]Result, 0)

	for _, result := range searchResults.Results {
		results = append(results, Result{result.PrefLabel, result.Uri})
	}

	return results, nil
}
