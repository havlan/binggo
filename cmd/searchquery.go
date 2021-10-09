package cmd

import "time"

type SearchQuery struct {
	Where string `schema:"where"`
	Query string `schema:"q"`
	Count int `schema:"count"`
	Offset int `schema:"offset"`
	When time.Time `schema:"when"`
	Safe bool `schema:"safe"`
}