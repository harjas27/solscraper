package query

type Query struct {
	Program string       `json:"program"`
	Filters []Match      `json:"filter"`
	Export  []ExportFunc `json:"export"`
}

type ExportFunc struct {
	Name     string `json:"name"`
	Function string `json:"value"`
}

type Match struct {
	Value      string
	MatchType  string
	MatchValue string
}
