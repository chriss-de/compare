package compare

import (
	"reflect"
	"time"
)

func (c *Comparer) cmpTime(path []string, a, b reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, a, b); err != nil || changed {
		return err
	}

	// Marshal and unmarshal time type will lose accuracy. Using unix nano to compare time type.
	au := getAsAny(a).(time.Time).UnixNano()
	bu := getAsAny(b).(time.Time).UnixNano()

	if au != bu {
		c.changes.add(CHANGE, path, getAsAny(a), getAsAny(b))
	}

	return nil
}
