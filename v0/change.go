package go_compare

// Changes stores a list of changed items
type Changes []Change

type ChangeType string

const (
	// ADD represents when an element has been added
	ADD ChangeType = "add"
	// CHANGE represents when an element has been updated
	CHANGE ChangeType = "change"
	// REMOVE represents when an element has been removed
	REMOVE ChangeType = "remove"
)

// Change stores information about a changed item
type Change struct {
	Type   ChangeType `json:"type"`
	Path   []string   `json:"path"`
	From   any        `json:"from"`
	To     any        `json:"to"`
	parent any        `json:"parent"`
}

func (c *Changes) add(t ChangeType, path []string, from any, to any, parent ...any) {
	change := Change{
		Type: t,
		Path: path,
		From: from,
		To:   to,
	}
	if len(parent) > 0 {
		change.parent = parent[0]
	}
	*c = append(*c, change)
}
