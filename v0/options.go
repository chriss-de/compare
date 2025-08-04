package compare

type OptsFunc func(d *Comparer) error

// WithSliceOrdering determines whether the ordering of items in a slice results in a change
func WithSliceOrdering(enabled bool) func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.sliceOrdering = enabled
		return nil
	}
}

// WithTagName sets the tag name to use when getting field names and options
func WithTagName(tag string) func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.tagName = tag
		return nil
	}
}

func WithCombinedIdentifierJoinString(joinSep rune) func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.combinedIdentifierJoinSep = joinSep
		return nil
	}
}

func WithSummarizeMissingStructs() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.summarizeMissingStructs = true
		return nil
	}
}

func WithStructMapKeys() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.structMapKeys = true
		return nil
	}
}

func WithEmbeddedStructsAsField() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.config.embeddedStructFields = false
		return nil
	}
}

//func WithStructIdInSlices() func(c *Comparer) error {
//	return func(c *Comparer) error {
//		c.useStructIdInSlices = true
//		return nil
//	}
//}
