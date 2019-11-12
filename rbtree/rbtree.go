package rbtree

import (
	"fmt"
	"github.com/shengmingzhu/datastructures/pair"
	"strings"
)

type RbTree struct {
	len  int
	root *node
	cmp  CmpFunc // cmp(key1, key2). It returns 0 if key1 == key2, returns 1 if key1 > key2, returns -1 if key1 < key2.
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

func New(f CmpFunc) *RbTree {
	return &RbTree{len: 0, root: nilNode, cmp: f}
}

func (t *RbTree) Len() int {
	return t.len
}

func (t *RbTree) IsEmpty() bool {
	return t.len == 0
}

// Search returns the value to key, or nil if not found.
// For example: if value := t.Search(key); value != nil { value found }
// O(logN)
func (t *RbTree) Get(key interface{}) (value interface{}) {
	p := t.search(key)
	if p == nilNode {
		return nil
	} else {
		return p.value
	}
}

// Put stores the key-value pair into RbTree.
// 1. If there is already a same key in RbTree, it will replace the value.
// 2. Otherwise, it will insert a new node with the key-value.
// O(logN)
func (t *RbTree) Put(key interface{}, value interface{}) {
	y := nilNode
	x := t.root
	for x != nilNode {
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
	z := newNodeForInsert(key, value, y)
	if y == nilNode {
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
func (t *RbTree) Delete(key interface{}) {
	z := t.search(key)
	if z == nilNode {
		return // not found
	}

	t.delete(z)
}

// Min returns the key-value to the minimum key, or nil if the tree is empty.
// For example: if key, value := t.Min(key); key != nil { found }
// O(logN)
func (t *RbTree) Min() (key, value interface{}) {
	p := t.root.min()
	if p == nilNode {
		return nil, nil
	} else {
		return p.key, p.value
	}
}

// Max returns the key-value to the maximum key, or nil if the tree is empty.
// For example: if key, value := t.Max(); key != nil { found }
// O(logN)
func (t *RbTree) Max() (key, value interface{}) {
	p := t.root.max()
	if p == nilNode {
		return nil, nil
	} else {
		return p.key, p.value
	}
}

// RangeAll traversals in ASC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) RangeAll() []*pair.Pair {
	return t.root.rangeAllAsc()
}

// RangeAllDesc traversals in DESC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) RangeAllDesc() []*pair.Pair {
	return t.root.rangeAllDesc()
}

// Range traversals in [minKey, maxKey] in ASC
// MinKey & MaxKey are all closed interval.
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) Range(minKey, maxKey interface{}) []*pair.Pair {
	return t.root.rangeAsc(minKey, maxKey, t.cmp)
}

// RangeN get num key-values which >= key in ASC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) RangeN(num int, key interface{}) []*pair.Pair {
	return t.root.rangeAscN(num, nil, key, t.cmp)
}

// RangeDescN get num key-values which <= key in DESC
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) RangeDescN(num int, key interface{}) []*pair.Pair {
	return t.root.rangeDescN(num, nil, key, t.cmp)
}

// RangeDesc traversals in [minKey, maxKey] in DESC
// MinKey & MaxKey are all closed interval.
// Pair.First: Key, Pair.Second: Value
// O(N)
func (t *RbTree) RangeDesc(minKey, maxKey interface{}) []*pair.Pair {
	return t.root.rangeDesc(minKey, maxKey, t.cmp)
}

// PopMin will delete the min node and return it.
// O(logN)
func (t *RbTree) PopMin() (key, value interface{}) {
	p := t.root.min()
	t.delete(p)
	return p.key, p.value
}

// PopMax will delete the max node and return it.
// O(logN)
func (t *RbTree) PopMax() (key, value interface{}) {
	p := t.root.max()
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
func (t *RbTree) String() string {
	step := t.getKeyMaxLen()
	step += 3 // one step [key]colour, example: [123456]B
	depth := t.root.getDepth()
	lDepth := t.root.getLeftDepth()
	lMove := uint(0)
	if depth > lDepth {
		lMove = step<<(depth-lDepth) - step
	}
	buffs := make([]*strings.Builder, depth)
	for i := uint(0); i < depth; i++ {
		buffs[i] = &strings.Builder{}
	}
	t.root.makeString(buffs, step, lMove, depth, 1, true, false)
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

func (t *RbTree) search(key interface{}) *node {
	p := t.root
	for p != nilNode {
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
func (t *RbTree) delete(z *node) {
	if z == nilNode {
		return
	}
	y := z
	yOriginalColor := y.color
	var x *node
	if z.left == nilNode {
		x = z.right
		t.transplant(z, z.right)
	} else if z.right == nilNode {
		x = z.left
		t.transplant(z, z.left)
	} else {
		y = z.right.min()
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
func (t *RbTree) leftRotate(x *node) {
	y := x.right

	x.right = y.left
	if y.left != nilNode {
		y.left.parent = x
	}

	y.parent = x.parent
	if x.parent == nilNode {
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
func (t *RbTree) rightRotate(x *node) {
	y := x.left

	x.left = y.right
	if y.right != nilNode {
		y.right.parent = x
	}

	y.parent = x.parent
	if x.parent == nilNode {
		t.root = y
	} else if x == x.parent.left {
		x.parent.left = y
	} else {
		x.parent.right = y
	}

	y.right = x
	x.parent = y
}

func (t *RbTree) transplant(u, v *node) {
	if u.parent == nilNode {
		t.root = v
	} else if u == u.parent.left {
		u.parent.left = v
	} else {
		u.parent.right = v
	}

	v.parent = u.parent
}

// O(logN)
func (t *RbTree) fixupInsert(z *node) {
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
func (t *RbTree) fixupDelete(x *node) {
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
func (t *RbTree) getKeyMaxLen() uint {
	minNode := t.root.min()
	lenMin := uint(len(fmt.Sprintf("%v", minNode.key)))
	maxNode := t.root.max()
	lenMax := uint(len(fmt.Sprintf("%v", maxNode.key)))
	var step uint
	if lenMin > lenMax {
		step = lenMin
	} else {
		step = lenMax
	}
	return step
}
