package rbtree

import (
	"fmt"
	"github.com/shengmingzhu/datastructures/pair"
	"strings"
)

type rbTree struct {
	len  int
	root *node
	cmp  CmpFunc // cmp(key1, key2). It returns 0 if key1 == key2, returns 1 if key1 > key2, returns -1 if key1 < key2.
	nil  *node
}

// CmpFunc such as CmpFunc(key1, key2).
// It returns 0 if key1 == key2, returns a number greater than 0 if key1 > key2, or less than 0 if key1 < key2.
/*
Example:
    if key1 == key2 { return 0 }
    else if key1 > key2 { return 1 }
    else { return -1 }
*/
type CmpFunc func(interface{}, interface{}) int

func New(f CmpFunc) *rbTree {
	nilNode := &node{color: black}
	return &rbTree{len: 0, root: nilNode, cmp: f, nil: nilNode}
}

func (t *rbTree) Len() int {
	return t.len
}

func (t *rbTree) IsEmpty() bool {
	return t.len == 0
}

// Search returns the value to key, or nil if not found.
// For example: if value, ok := t.Search(key); ok { value found }
// O(logN)
func (t *rbTree) Get(key interface{}) (value interface{}, ok bool) {
	p := t.search(key)
	if p == t.nil {
		return nil, false
	} else {
		return p.value, true
	}
}

// Put stores the key-value pair into rbTree.
// 1. If there is already a same key in rbTree, it will replace the value.
// 2. Otherwise, it will insert a new node with the key-value.
// O(logN)
func (t *rbTree) Put(key interface{}, value interface{}) {
	y := t.nil
	x := t.root
	for x != t.nil {
		y = x
		if t.cmp(x.key, key) == 0 {
			x.value = value
			return // if found, save value and return
		} else if t.cmp(x.key, key) > 0 {
			x = x.left
		} else {
			x = x.right
		}
	}

	// if not found, we insert a new node
	z := t.newNodeForInsert(key, value, y)
	if y == t.nil {
		t.root = z
	} else if t.cmp(y.key, key) > 0 {
		y.left = z
	} else {
		y.right = z
	}
	t.len++

	t.fixupInsert(z)
}

// O(logN)
func (t *rbTree) Delete(key interface{}) {
	z := t.search(key)
	if z == t.nil {
		return // not found
	}

	t.delete(z)
}

// Min returns the key-value to the minimum key, or nil if the tree is empty.
// For example: if key, value := t.Min(key); key != nil { found }
// O(logN)
func (t *rbTree) Min() (key, value interface{}) {
	p := t.min(t.root)
	if p == t.nil {
		return nil, nil
	} else {
		return p.key, p.value
	}
}

// Max returns the key-value to the maximum key, or nil if the tree is empty.
// For example: if key, value := t.Max(); key != nil { found }
// O(logN)
func (t *rbTree) Max() (key, value interface{}) {
	p := t.max(t.root)
	if p == t.nil {
		return nil, nil
	} else {
		return p.key, p.value
	}
}

// Keys traversals in ASC
// O(N)
func (t *rbTree) Keys() []interface{} {
	pos := 0
	res := make([]interface{}, t.len)
	t.rangeKeysAsc(t.root, res, &pos)
	return res
}

// Values traversals in ASC
// O(N)
func (t *rbTree) Values() []interface{} {
	pos := 0
	res := make([]interface{}, t.len)
	t.rangeValuesAsc(t.root, res, &pos)
	return res
}

// RangeAll traversals in ASC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) RangeAll() []pair.Pair {
	pos := 0
	res := make([]pair.Pair, t.len)
	t.rangeAllAsc(t.root, res, &pos)
	return res
}

// RangeAllDesc traversals in DESC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) RangeAllDesc() []pair.Pair {
	pos := 0
	res := make([]pair.Pair, t.len)
	t.rangeAllDesc(t.root, res, &pos)
	return res
}

// Range traversals in [minKey, maxKey] in ASC
// MinKey & MaxKey are all closed interval.
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) Range(minKey, maxKey interface{}) []pair.Pair {
	return t.rangeAsc(t.root, nil, minKey, maxKey, t.cmp)
}

// RangeN get num key-values which >= key in ASC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) RangeN(num int, key interface{}) []pair.Pair {
	return t.rangeAscN(t.root, nil, num, key, t.cmp)
}

// RangeDescN get num key-values which <= key in DESC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) RangeDescN(num int, key interface{}) []pair.Pair {
	return t.rangeDescN(t.root, nil, num, key, t.cmp)
}

// RangeDesc traversals in [minKey, maxKey] in DESC
// MinKey & MaxKey are all closed interval.
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *rbTree) RangeDesc(minKey, maxKey interface{}) []pair.Pair {
	return t.rangeDesc(t.root, nil, minKey, maxKey, t.cmp)
}

// PopMin will delete the min node and return it.
// O(logN)
func (t *rbTree) PopMin() (key, value interface{}) {
	p := t.min(t.root)
	t.delete(p)
	return p.key, p.value
}

// PopMax will delete the max node and return it.
// O(logN)
func (t *rbTree) PopMax() (key, value interface{}) {
	p := t.max(t.root)
	t.delete(p)
	return p.key, p.value
}

// String is very useful when debugging
// Example: fmt.Println(t) will print as follows:
/*
                              [ 6]B
          [ 3]B                                   [13]R
[ 2]R               [ 5]R               [ 8]B               [15]B
                                             [10]R
*/
// Deprecated: only for debugging, unstable function
func (t *rbTree) String() string {
	step := t.getKeyMaxLen()
	step += 3 // one step [key]colour, example: [123456]B
	depth := t.getDepth(t.root)
	lDepth := t.getLeftDepth(t.root)
	lMove := uint(0)
	if depth > lDepth {
		lMove = step<<(depth-lDepth) - step
	}
	buffs := make([]*strings.Builder, depth)
	for i := uint(0); i < depth; i++ {
		buffs[i] = &strings.Builder{}
	}
	t.makeString(t.root, buffs, step, lMove, depth, 1, true, false)
	/*check, ok := t.root.check()
	strCheck := fmt.Sprintf("Tree check is %v, MaxH %v, MinH %v, BlackH %v, Len %v, Root colour: %v", ok, check.MaxH, check.MinH, check.BlackH, check.Len, check.Colour)
	ss := []string{strCheck}*/
	var ss []string
	for i := 0; i < len(buffs); i++ {
		ss = append(ss, buffs[i].String())
	}
	str := strings.Join(ss, fmt.Sprintln(""))

	return str
}

func (t *rbTree) search(key interface{}) *node {
	p := t.root
	for p != t.nil {
		cmp := t.cmp(p.key, key)
		if cmp == 0 {
			break
		} else if cmp > 0 {
			p = p.left
		} else {
			p = p.right
		}
	}
	return p
}

// O(logN)
func (t *rbTree) delete(z *node) {
	if z == t.nil {
		return
	}
	y := z
	yOriginalColor := y.color
	var x *node
	if z.left == t.nil {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == t.nil {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = t.min(z.right)
		yOriginalColor = y.color
		x = y.right
		if y.parent == z {
			// This line is useful.
			// If x == nilNode, x is still possible to pass to t.fixupDelete(x),
			// and t.fixupDelete(x) need x.parent is valid.
			x.parent = y
		} else {
			t.transplant(y, y.right)
			y.right = z.right
			y.right.parent = y
		}
		t.transplant(z, y)
		y.left = z.left
		y.left.parent = y
		y.color = z.color
	}
	t.len--

	if yOriginalColor == black {
		t.fixupDelete(x)
	}
}

// O(1)
func (t *rbTree) leftRotate(x *node) {
	y := x.right

	x.right = y.left
	if y.left != t.nil {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == t.nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.left = x
	x.parent = y
}

// O(1)
func (t *rbTree) rightRotate(x *node) {
	y := x.left

	x.left = y.right
	if y.right != t.nil {
		y.right.parent = x
	}

	y.parent = x.parent
	if x.parent == t.nil {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}

func (t *rbTree) transplant(u, v *node) {
	if u.parent == t.nil {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	v.parent = u.parent
}

// O(logN)
func (t *rbTree) fixupInsert(z *node) {
	// The necessary conditions for each entry into the loop:
	// 1. z.color == red
	// 2. if t.root == z.parent, z.parent must be black
	// 3. z != t.root and z.parent must be red
	for z.parent.color == red { // if z == t.root, z.parent must be nilNode, and nilNode.color == black
		if z.parent == z.parent.parent.left {
			y := z.parent.parent.right // uncle node
			if y.color == red {        // case 1
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.right { // case 2
					z = z.parent
					t.leftRotate(z)
				}

				// case 3
				z.parent.color = black
				z.parent.parent.color = red
				t.rightRotate(z.parent.parent)
			}
		} else { // z.parent == z.parent.parent.right
			y := z.parent.parent.left // uncle node
			if y.color == red {       // case 1
				z.parent.color = black
				y.color = black
				z.parent.parent.color = red
				z = z.parent.parent
			} else {
				if z == z.parent.left { // case 2
					z = z.parent
					t.rightRotate(z)
				}

				// case 3
				z.parent.color = black
				z.parent.parent.color = red
				t.leftRotate(z.parent.parent)
			}
		}
	}

	t.root.color = black // deal with z == t.root, case z must be red
}

// O(logN)
func (t *rbTree) fixupDelete(x *node) {
	// x is a node carries extra black, it can be red-black or black-black.
	//   1. if x == t.root, we can just remove the extra black.
	//   2. if x.color == red, we can change x to black.
	// In other cases, we fixup in loop.
	for x != t.root && x.color == black {
		if x == x.parent.left {
			w := x.parent.right
			if w.color == red {
				w.color = black        // case 1
				x.parent.color = red   // case 1
				t.leftRotate(x.parent) // case 1
				w = x.parent.right     // case 1
			}
			if w.left.color == black && w.right.color == black {
				w.color = red // case 2
				x = x.parent  // case 2
			} else {
				if w.right.color == black {
					w.left.color = black // case 3
					w.color = red        // case 3
					t.rightRotate(w)     // case 3
					w = x.parent.right   // case 3
				}
				w.color = x.parent.color // case 4
				x.parent.color = black   // case 4
				w.right.color = black    // case 4
				t.leftRotate(x.parent)   // case 4
				x = t.root               // case 4
			}
		} else {
			w := x.parent.left
			if w.color == red {
				w.color = black         // case 1
				x.parent.color = red    // case 1
				t.rightRotate(x.parent) // case 1
				w = x.parent.left       // case 1
			}

			// w.color == black
			if w.left.color == black && w.right.color == black {
				w.color = red // case 2
				x = x.parent  // case 2
			} else {
				if w.left.color == black {
					w.right.color = black // case 3
					w.color = red         // case 3
					t.leftRotate(w)       // case 3
					w = x.parent.left     // case 3
				}
				w.color = x.parent.color // case 4
				x.parent.color = black   // case 4
				w.left.color = black     // case 4
				t.rightRotate(x.parent)  // case 4
				x = t.root               // case 4
			}
		}
	}
	x.color = black
}

// O(logN)
func (t *rbTree) getKeyMaxLen() uint {
	minNode := t.min(t.root)
	lenMin := uint(len(fmt.Sprintf("%v", minNode.key)))
	maxNode := t.max(t.root)
	lenMax := uint(len(fmt.Sprintf("%v", maxNode.key)))
	var step uint
	if lenMin > lenMax {
		step = lenMin
	} else {
		step = lenMax
	}
	return step
}

type node struct {
	key    interface{}
	value  interface{}
	parent *node // parent
	left   *node // left child
	right  *node // right child
	color  colours
}

type colours uint8

const (
	red colours = iota
	black
)

// newNodeForInsert returns a pointer to the new node containing the key/value, the new node must be red
func (t *rbTree) newNodeForInsert(key interface{}, value interface{}, parent *node) *node {
	return &node{key: key, value: value, color: red, parent: parent, left: t.nil, right: t.nil}
}

// O(logN)
func (t *rbTree) min(n *node) *node {
	p := n
	for p != t.nil && p.left != t.nil {
		p = p.left
	}
	return p
}

// O(logN)
func (t *rbTree) max(n *node) *node {
	p := n
	for p != t.nil && p.right != t.nil {
		p = p.right
	}
	return p
}

// O(logN)
func (t *rbTree) successor(n *node) *node {
	if n.right != t.nil {
		return t.min(n.right)
	}

	x := n
	y := x.parent
	for y != t.nil && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

// O(logN)
func (t *rbTree) predecessor(n *node) *node {
	if n.left != t.nil {
		return t.max(n.left)
	}

	return n.parent
}

func (t *rbTree) rangeAllAsc(n *node, res []pair.Pair, pos *int) {
	if n == t.nil {
		return
	}

	if n.left != t.nil {
		t.rangeAllAsc(n.left, res, pos)
	}
	res[*pos].First = n.key
	res[*pos].Second = n.value
	*pos++
	if n.right != t.nil {
		t.rangeAllAsc(n.right, res, pos)
	}
}

func (t *rbTree) rangeKeysAsc(n *node, res []interface{}, pos *int) {
	if n == t.nil {
		return
	}
	if n.left != t.nil {
		t.rangeKeysAsc(n.left, res, pos)
	}
	res[*pos] = n.key
	*pos++
	if n.right != t.nil {
		t.rangeKeysAsc(n.right, res, pos)
	}
}

func (t *rbTree) rangeValuesAsc(n *node, res []interface{}, pos *int) {
	if n == t.nil {
		return
	}
	if n.left != t.nil {
		t.rangeValuesAsc(n.left, res, pos)
	}
	res[*pos] = n.value
	*pos++
	if n.right != t.nil {
		t.rangeValuesAsc(n.right, res, pos)
	}
}

func (t *rbTree) rangeAllDesc(n *node, res []pair.Pair, pos *int) {
	if n == t.nil {
		return
	}

	if n.right != t.nil {
		t.rangeAllDesc(n.right, res, pos)
	}
	res[*pos].First = n.key
	res[*pos].Second = n.value
	*pos++
	if n.left != t.nil {
		t.rangeAllDesc(n.left, res, pos)
	}
}

func (t *rbTree) rangeAsc(n *node, res []pair.Pair, minKey, maxKey interface{}, cmp CmpFunc) []pair.Pair {
	if n == t.nil {
		return res
	}

	cmpMin, cmpMax := cmp(n.key, minKey), cmp(n.key, maxKey) // cmp() may takes some time, so we just cmp one time.
	if cmpMin > 0 {
		res = t.rangeAsc(n.left, res, minKey, maxKey, cmp)
	}
	if cmpMin >= 0 && cmpMax <= 0 {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if cmpMax < 0 {
		res = t.rangeAsc(n.right, res, minKey, maxKey, cmp)
	}
	return res
}

func (t *rbTree) rangeAscN(n *node, res []pair.Pair, num int, key interface{}, cmp CmpFunc) []pair.Pair {
	if n == t.nil {
		return res
	}

	iCmp := cmp(n.key, key) // cmp() may takes some time, so we just cmp one time.
	if iCmp > 0 && len(res) < num {
		res = t.rangeAscN(n.left, res, num, key, cmp)
	}
	if iCmp >= 0 && len(res) < num {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if len(res) < num {
		res = t.rangeAscN(n.right, res, num, key, cmp)
	}
	return res
}

func (t *rbTree) rangeDesc(n *node, res []pair.Pair, minKey, maxKey interface{}, cmp CmpFunc) []pair.Pair {
	if n == t.nil {
		return res
	}

	cmpMin, cmpMax := cmp(n.key, minKey), cmp(n.key, maxKey) // cmp() may takes some time, so we just cmp one time.
	if cmpMax < 0 {
		res = t.rangeDesc(n.right, res, minKey, maxKey, cmp)
	}
	if cmpMin >= 0 && cmpMax <= 0 {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if cmpMin > 0 {
		res = t.rangeDesc(n.left, res, minKey, maxKey, cmp)
	}
	return res
}

func (t *rbTree) rangeDescN(n *node, res []pair.Pair, num int, key interface{}, cmp CmpFunc) []pair.Pair {
	if n == t.nil {
		return res
	}

	iCmp := cmp(n.key, key) // cmp() may takes some time, so we just cmp one time.
	if iCmp < 0 && len(res) < num {
		res = t.rangeDescN(n.right, res, num, key, cmp)
	}
	if iCmp <= 0 && len(res) < num {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if len(res) < num {
		res = t.rangeDescN(n.left, res, num, key, cmp)
	}
	return res
}

type Check struct {
	MaxH   int     `json:"maxH"`
	MinH   int     `json:"minH"`
	BlackH int     `json:"blackH"`
	Len    int     `json:"len"`
	Colour colours `json:"colour"`
}

// check checks if n is a rbTree
// if ok == true, the tree is a rbTree.
func (t *rbTree) check(n *node) (c *Check, ok bool) {
	if n == t.nil {
		return &Check{Colour: black}, true
	}
	c = &Check{MaxH: 1, MinH: 1, Len: 1, Colour: n.color}
	if n.color == black {
		c.BlackH = 1
	}
	if n.left == t.nil && n.right == t.nil {
		return c, true
	} else if n.left != t.nil && n.right != t.nil {
		var cl, cr *Check
		ok := false
		if cl, ok = t.check(n.left); !ok {
			return c, false
		}
		if cr, ok = t.check(n.right); !ok {
			return c, false
		}

		if (c.Colour == red && (cl.Colour == red || cr.Colour == red)) || cl.BlackH != cr.BlackH {
			return c, false
		}

		if cl.MaxH > cr.MaxH {
			c.MaxH += cl.MaxH
		} else {
			c.MaxH += cr.MaxH
		}

		if cl.MinH < cr.MinH {
			c.MinH += cl.MinH
		} else {
			c.MinH += cr.MinH
		}

		c.BlackH += cl.BlackH
		c.Len += cl.Len + cr.Len
	} else if n.right == t.nil {
		cl, ok := t.check(n.left)
		if !ok || cl.Colour != red || cl.Len != 1 {
			return c, false
		}

		c.Len += cl.Len
		c.BlackH += cl.BlackH
		c.MinH += cl.MinH
		c.MaxH += cl.MaxH
	} else {
		cr, ok := t.check(n.right)
		if !ok || cr.Colour != red || cr.Len != 1 {
			return c, false
		}

		c.Len += cr.Len
		c.BlackH += cr.BlackH
		c.MinH += cr.MinH
		c.MaxH += cr.MaxH
	}

	return c, true
}

// O(N)
func (t *rbTree) getDepth(n *node) uint {
	if n == t.nil {
		return 0
	}
	ld := t.getDepth(n.left)
	rd := t.getDepth(n.right)

	if ld > rd {
		return ld + 1
	} else {
		return rd + 1
	}
}

// O(logN)
func (t *rbTree) getLeftDepth(n *node) uint {
	if n == t.nil {
		return 0
	}
	return t.getLeftDepth(n.left) + 1
}

// Deprecated: only for debugging, unstable function
func (t *rbTree) makeString(n *node, buffs []*strings.Builder, step, lMove, tDepth, nDepth uint, ifRowFirst, ifParentRowFirst bool) {
	if n == t.nil {
		return
	}
	space := (step << (tDepth - nDepth)) - step
	if ifRowFirst {
		if space >= lMove {
			space -= lMove
		}
	} else {
		space = (space << 1) + step
		if ifParentRowFirst && n.parent.left == t.nil {
			space = step * (((space / step) >> 1) + 1)
		}
	}
	for i := space; i > 0; i-- {
		buffs[nDepth-1].WriteString(" ")
	}
	// write key "[123456]B"
	buffs[nDepth-1].WriteString("[")
	strKey := fmt.Sprint(n.key)
	for i := int(step) - 3 - len(strKey); i > 0; i-- {
		buffs[nDepth-1].WriteString(" ")
	}
	buffs[nDepth-1].WriteString(strKey)
	if n.color == red {
		buffs[nDepth-1].WriteString("]R")
	} else {
		buffs[nDepth-1].WriteString("]B")
	}

	if n.left != t.nil {
		t.makeString(n.left, buffs, step, lMove, tDepth, nDepth+1, ifRowFirst, false)
	} else if !ifRowFirst {
		for i := nDepth + 1; i <= tDepth; i++ {
			rStep := spaceCount(step, tDepth, i)
			rSpace := rStep << (i - nDepth - 1)
			for j := rSpace; j > 0; j-- {
				buffs[i-1].WriteString(" ")
			}
		}
	}

	if n.right != t.nil {
		if n.left == t.nil {
			t.makeString(n.right, buffs, step, lMove, tDepth, nDepth+1, false, true)
		} else {
			t.makeString(n.right, buffs, step, lMove, tDepth, nDepth+1, false, false)
		}
	} else {
		if ifRowFirst && n.left == t.nil {
			rSpace := step << (tDepth - nDepth)
			for i := tDepth; i > nDepth; i-- {
				if i != tDepth {
					rSpace -= step << (tDepth - i - 1)
				}
				for j := rSpace; j > 0; j-- {
					buffs[i-1].WriteString(" ")
				}
			}
		} else {
			for i := nDepth + 1; i <= tDepth; i++ {
				rStep := spaceCount(step, tDepth, i)
				rSpace := rStep << (i - nDepth - 1)
				for j := rSpace; j > 0; j-- {
					buffs[i-1].WriteString(" ")
				}
			}
		}
	}
}

func spaceCount(step, tDepth, nDepth uint) uint {
	return step << (tDepth - nDepth + 1)
}
