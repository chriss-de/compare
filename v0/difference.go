package compare

import (
	"iter"
)

// Differences stores a list of changed items
type Differences []Difference

type DiffType string

type DiffFilterFunc func(Difference) bool

const (
	// ADD represents when an element has been added
	ADD DiffType = "add"
	// CHANGE represents when an element has been updated
	CHANGE DiffType = "change"
	// REMOVE represents when an element has been removed
	REMOVE DiffType = "remove"
)

// Difference stores information about a changed item
type Difference struct {
	Type   DiffType `json:"type"`
	Path   []string `json:"path"`
	From   any      `json:"from"`
	To     any      `json:"to"`
	parent any      `json:"parent"`
}

func (d *Differences) add(t DiffType, path []string, from any, to any, parent ...any) {
	diff := Difference{
		Type: t,
		Path: path,
		From: from,
		To:   to,
	}
	if len(parent) > 0 {
		diff.parent = parent[0]
	}
	*d = append(*d, diff)
}

func (d *Differences) GetDifferences(filterFunc ...DiffFilterFunc) iter.Seq[Difference] {
	return func(yield func(diff Difference) bool) {
		for _, k := range *d {
			for _, ff := range filterFunc {
				// filterFunc ar AND
				var yieldIt = true
				if !ff(k) {
					yieldIt = false
				}
				if yieldIt && !yield(k) {
					return
				}
			}
		}
	}
}
