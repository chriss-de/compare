package compare

type CompareOptsFunc func(d *Comparer) error

// WithSliceOrdering determines whether the ordering of items in a slice results in a change
func WithSliceOrdering(enabled bool) func(c *Comparer) error {
	return func(c *Comparer) error {
		c.sliceOrdering = enabled
		return nil
	}
}

// WithTagName sets the tag name to use when getting field names and options
func WithTagName(tag string) func(c *Comparer) error {
	return func(c *Comparer) error {
		c.tagName = tag
		return nil
	}
}

func WithSummarizeMissingStructs() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.summarizeMissingStructs = true
		return nil
	}
}

func WithStructMapKeys() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.structMapKeys = true
		return nil
	}
}

func WithEmbeddedStructsAsField() func(c *Comparer) error {
	return func(c *Comparer) error {
		c.embeddedStructFields = false
		return nil
	}
}

//func WithStructIdInSlices() func(c *Comparer) error {
//	return func(c *Comparer) error {
//		c.useStructIdInSlices = true
//		return nil
//	}
//}
