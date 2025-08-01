package compare

import (
	"reflect"
	"time"
)

func (c *Comparer) cmpTime(path []string, left, right reflect.Value, parent any) error {
	if changed, err := c.cmpDefault(path, left, right); err != nil || changed {
		return err
	}

	// Marshal and unmarshal time type will lose accuracy. Using unix nano to compare time type.
	au := getAsAny(left).(time.Time).UnixNano()
	bu := getAsAny(right).(time.Time).UnixNano()

	if au != bu {
		c.changes.add(CHANGE, path, getAsAny(left), getAsAny(right))
	}

	return nil
}
