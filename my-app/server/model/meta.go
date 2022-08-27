package model

import "ramdeuter.org/solscraper/query"

type Meta struct {
	Name  string      `json:"name"`
	Query query.Query `json:"query"`
}
