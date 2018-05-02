package agrovoc

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

type SearchResults struct {
	Uri string `json:"uri"`
	Results []SearchResult `json:"results"`
}

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

type Result struct {
	Label string
	Uri string
}