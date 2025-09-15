package errors

import (
	"reflect"
	"sort"
)

type Empty struct{}

type String map[string]Empty

func NewString(items ...string) String {
	new := make(String)
	new.Insert(items...)

	return new
}

func MapKeysToString(theMap any) String {
	new := NewString()

	reflectedMap := reflect.ValueOf(theMap)

	for _, v := range reflectedMap.MapKeys() {
		new.Insert(v.Interface().(string))
	}

	return new
}

func (s String) Insert(items ...string) {
	for _, item := range items {
		s[item] = Empty{}
	}
}

func (s String) Delete(item string) {
	delete(s, item)
}

func (s String) Has(item string) bool {
	_, ok := s[item]
	return ok
}

func (s String) HasAny(items ...string) bool {
	for _, item := range items {
		if s.Has(item) {
			return true
		}
	}

	return false
}

func (s String) HasAll(items ...string) bool {
	for _, item := range items {
		if !s.Has(item) {
			return false
		}
	}

	return true
}

func (s String) Difference(b String) String {
	d := NewString()

	for key := range s {
		if !b.Has(key) {
			d.Insert(key)
		}
	}

	return d
}

func (s String) InterSection(b String) String {
	i := NewString()
	var walk String
	var other String

	if len(s) < len(b) {
		walk = s
		other = b
	} else {
		walk = b
		other = s
	}

	for key := range walk {
		if other.Has(key) {
			i.Insert(key)
		}
	}

	return i
}

func (s String) Union(b String) String {
	u := NewString()

	for key := range s {
		u.Insert(key)
	}

	for key := range b {
		u.Insert(key)
	}

	return u
}

func (s String) IsSuperSet(b String) bool {
	for key := range b {
		if !s.Has(key) {
			return false
		}
	}

	return true
}

func (s String) IsEqual(b String) bool {
	return len(s) == len(b) && s.IsSuperSet(b)
}

type SortableStringList []string

func (l SortableStringList) Len() int { return len(l) }

func (l SortableStringList) Less(i, j int) bool { return l[i] < l[j] }

func (l SortableStringList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }

func (s String) ConvertToSortedList() SortableStringList {
	l := make(SortableStringList, 0, len(s))

	for key := range s {
		l = append(l, key)
	}

	sort.Sort(l)

	return l
}

func (s String) ConvertToUnsortedList() SortableStringList {
	l := make(SortableStringList, 0, len(s))

	for key := range s {
		l = append(l, key)
	}

	return l
}

func (s String) PopAny() (string, bool) {
	for key := range s {
		s.Delete(key)
		return key, true
	}

	var ZeroValue string

	return ZeroValue, false
}

func (s String) Len() int {
	return len(s)
}
