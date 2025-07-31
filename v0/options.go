package go_compare

// WithSliceOrdering determines whether the ordering of items in a slice results in a change
func WithSliceOrdering(enabled bool) func(d *Comparer) error {
	return func(d *Comparer) error {
		d.sliceOrdering = enabled
		return nil
	}
}

// WithTagName sets the tag name to use when getting field names and options
func WithTagName(tag string) func(d *Comparer) error {
	return func(d *Comparer) error {
		d.tagName = tag
		return nil
	}
}

func WithSummarizeMissingStructs() func(d *Comparer) error {
	return func(d *Comparer) error {
		d.summarizeMissingStructs = true
		return nil
	}
}

func WithStructMapKeys() func(d *Comparer) error {
	return func(d *Comparer) error {
		d.structMapKeys = true
		return nil
	}
}

func WithEmbeddedStructsAsField() func(d *Comparer) error {
	return func(d *Comparer) error {
		d.embeddedStructFields = false
		return nil
	}
}

//func WithStructIdInSlices() func(d *Comparer) error {
//	return func(d *Comparer) error {
//		d.useStructIdInSlices = true
//		return nil
//	}
//}
