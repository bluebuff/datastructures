/*
Key words prompt, with weight.
*/
package tirekwp

import (
	"fmt"
	"github.com/shengmingzhu/orderedmap"
	"strings"
	"unsafe"
)

type TireKWP struct {
	root         *node
	maxSortedLen int
	len          int // count of keywords
	nNodes       int // count of nodes
}

type Keyword struct {
	Weight int
	Str    string
	str    []rune
}

func New(maxSortedLen int) *TireKWP {
	return &TireKWP{
		root:         newNode(nil),
		maxSortedLen: maxSortedLen,
		len:          0,
		nNodes:       1,
	}
}

func newNode(key *Keyword) *node {
	n := &node{
		key:    key,
		next:   make(map[rune]*node),
		sorted: orderedmap.NewAny(cmp),
	}

	if key != nil {
		n.sorted.Put(key, nil)
	}

	return n
}

func (t *TireKWP) Put(str string, weight int) {
	if len(str) <= 0 {
		panic("Can't put an empty string to tireKWP.")
	}
	key := Keyword{str: []rune(str), Str: str, Weight: weight}
	if len(key.str) <= 0 {
		panic(fmt.Sprintf("We have a problem when converting string[%s] to rune.", str))
	}

	n := t.get(key.str)
	if n != nil && n.key != nil && len(n.key.Str) == len(str) {
		return
	}

	t.put(&key)
	t.len++
}

// key must not in t. If not sure, must get() and delete() key first.
func (t *TireKWP) put(key *Keyword) {
	// 1. pos = len has traversal
	// 2. pos point to the next rune
	pos := 0
	now, ok := t.root.next[key.str[pos]]
	if !ok {
		t.root.next[key.str[pos]] = newNode(key)
		t.nNodes++
		t.root.adjustSorted(key, t.maxSortedLen)
		return
	}
	t.root.adjustSorted(key, t.maxSortedLen)

	for {
		pos++
		// Case 1: key traversal completed, store key to now.
		if pos == len(key.str) {
			if now.key != nil {
				// Case 1.1: now is leaf node
				now.next[now.key.str[pos]] = newNode(now.key)
				t.nNodes++
			}
			// store key to now
			now.key = key
			// key's weight may changed, so we adjust now.sorted
			now.adjustSorted(key, t.maxSortedLen)
			break
		}

		// Case 2: now is leaf node
		if len(now.next) <= 0 {
			k2 := now.key
			// Handling same prefixes in loop
			for pos < len(key.str) && pos < len(k2.str) && key.str[pos] == k2.str[pos] {
				newN := newNode(nil)
				t.nNodes++
				now.next[key.str[pos]] = newN
				now.adjustSorted(key, t.maxSortedLen)
				now.adjustSorted(k2, t.maxSortedLen)
				now.key = nil
				now = newN
				pos++
			}
			if pos == len(key.str) { // Case 2.1: key traversal completed
				now.key = key
				now.adjustSorted(key, t.maxSortedLen)

				now.next[k2.str[pos]] = newNode(k2)
				t.nNodes++
				now.adjustSorted(k2, t.maxSortedLen)
			} else if pos == len(k2.str) { // Case 2.2: k2 traversal completed
				now.key = k2
				now.adjustSorted(k2, t.maxSortedLen)

				now.next[key.str[pos]] = newNode(key)
				t.nNodes++
				now.adjustSorted(key, t.maxSortedLen)
			} else { // Case 2.3: fork
				now.key = nil
				now.adjustSorted(key, t.maxSortedLen)
				now.adjustSorted(k2, t.maxSortedLen)

				now.next[key.str[pos]] = newNode(key)
				t.nNodes++
				now.adjustSorted(key, t.maxSortedLen)
				now.next[k2.str[pos]] = newNode(k2)
				t.nNodes++
				now.adjustSorted(k2, t.maxSortedLen)
			}
			break
		}

		// Case 3: now is not a leaf node
		next, ok := now.next[key.str[pos]]
		if !ok {
			now.next[key.str[pos]] = newNode(key)
			t.nNodes++
			now.adjustSorted(key, t.maxSortedLen)
			break
		}

		now.adjustSorted(key, t.maxSortedLen)
		now = next
	}
}

func (t *TireKWP) Get(str string) []string {
	var keys []interface{}
	if str == "" {
		keys = t.root.sorted.Keys()
	} else {
		r := []rune(str)
		n := t.get(r)
		if n != nil {
			keys = n.sorted.Keys()
		}
	}

	res := *((*[]string)(unsafe.Pointer(&keys)))
	for i := range keys {
		res[i] = keys[i].(*Keyword).Str
	}
	return res
}

func (t *TireKWP) GetKWs(str string) []*Keyword {
	var keys []interface{}
	if str == "" {
		keys = t.root.sorted.Keys()
	} else {
		r := []rune(str)
		n := t.get(r)
		if n != nil {
			keys = n.sorted.Keys()
		}
	}

	res := make([]*Keyword, len(keys))
	for i := range keys {
		res[i] = keys[i].(*Keyword)
	}
	return res
}

func (t *TireKWP) get(str []rune) *node {
	now := t.root
	ok := false
	for pos := 0; pos < len(str); pos++ {
		if len(now.next) <= 0 {
			if now.key != nil && len(now.key.str) >= len(str) {
				for i := pos; i < len(str); i++ {
					if str[i] != now.key.str[i] {
						return nil
					}
				}
				break // found
			} else {
				return nil
			}
		}
		now, ok = now.next[str[pos]]
		if !ok {
			return nil
		}
	}

	if now.key != nil {
		for i := range str {
			if str[i] != now.key.str[i] {
				return nil
			}
		}
		return now
	} else {
		return now
	}
}

func (t *TireKWP) Len() int {
	return t.len
}

func (t *TireKWP) Count() int {
	return t.nNodes
}

type node struct {
	/*
		1. If a word ends here, key will point to it.
		2. If there is only one word left, for save memory, we will not continue to allocate nodes, but directly point to the word with key.
		3. Other times, key == nil
	*/
	key    *Keyword
	next   map[rune]*node // For now, hash-map is fastest for search.
	sorted orderedmap.Any // The ordered keywords of each node are maintained during put() and delete(), so that get() can get quick response.
}

func (n *node) adjustSorted(key *Keyword, maxLen int) {
	n.sorted.Put(key, nil)
	if n.sorted.Len() > maxLen {
		_, _ = n.sorted.PopMax() // max is the least weight
	}
}

// cmp compare key1 and key2 for orderedmap
// Level 1, DESC of weight.
// Level 2, if weights are same, ASC of string
func cmp(key1, key2 interface{}) int {
	k1, k2 := key1.(*Keyword), key2.(*Keyword)
	c := strings.Compare(k1.Str, k2.Str)
	if c == 0 {
		return 0
	} else if k1.Weight != k2.Weight {
		return k2.Weight - k1.Weight // DESC of weight
	} else {
		return c // if weights are same, ASC of string
	}
}
