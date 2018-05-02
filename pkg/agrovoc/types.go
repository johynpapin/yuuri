package agrovoc

// SearchRequest represents a request to agrovoc.
type SearchRequest struct {
	Query string
	Lang string
	Labellang string
	Vocab string
	Type string
	Parent string
	Group string
	Maxhits int
	Offset int
	Fields []string
	Unique bool
}

// SearchResults represents the results of an agrovoc query.
type SearchResults struct {
	Uri string `json:"uri"`
	Results []SearchResult `json:"results"`
}

// SearchResult represents a result of an agrovoc query.
type SearchResult struct {
	Uri string `json:"uri"`
	Type []string `json:"type"`
	PrefLabel string `json:"prefLabel"`
	AltLabel string `json:"altLabel"`
	HiddenLabel string `json:"hiddenLabel"`
	Lang string `json:"lang"`
	Vocab string `json:"vocab"`
	Exvocab string `json:"exvocab"`
	Notation string `json:"notation"`
}

// Result represents a uri/label couple from agrovoc.
type Result struct {
	Label string
	Uri string
}