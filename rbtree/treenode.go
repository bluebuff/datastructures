package rbtree

import (
	"fmt"
	"github.com/shengmingzhu/datastructures/pair"
	"strings"
)

type node struct {
	key    interface{}
	value  interface{}
	parent *node // parent
	left   *node // left child
	right  *node // right child
	color  colours
}

var nilNode *node

type colours uint8

const (
	red colours = iota
	black
)

func init() {
	nilNode = &node{color: black}
}

// newNodeForInsert returns a pointer to the new node containing the key/value, the new node must be red
func newNodeForInsert(key interface{}, value interface{}, parent *node) *node {
	return &node{key: key, value: value, color: red, parent: parent, left: nilNode, right: nilNode}
}

// O(logN)
func (n *node) min() *node {
	p := n
	for p != nilNode && p.left != nilNode {
		p = p.left
	}
	return p
}

// O(logN)
func (n *node) max() *node {
	p := n
	for p != nilNode && p.right != nilNode {
		p = p.right
	}
	return p
}

// O(logN)
func (n *node) successor() *node {
	if n.right != nilNode {
		return n.right.min()
	}

	x := n
	y := x.parent
	for y != nilNode && x == y.right {
		x = y
		y = y.parent
	}
	return y
}

// O(logN)
func (n *node) predecessor() *node {
	if n.left != nilNode {
		return n.left.max()
	}

	return n.parent
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
func (n *node) check() (c *Check, ok bool) {
	if n == nilNode {
		return &Check{Colour: black}, true
	}
	c = &Check{MaxH: 1, MinH: 1, Len: 1, Colour: n.color}
	if n.color == black {
		c.BlackH = 1
	}
	if n.left == nilNode && n.right == nilNode {
		return c, true
	} else if n.left != nilNode && n.right != nilNode {
		var cl, cr *Check
		ok := false
		if cl, ok = n.left.check(); !ok {
			return c, false
		}
		if cr, ok = n.right.check(); !ok {
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
	} else if n.right == nilNode {
		cl, ok := n.left.check()
		if !ok || cl.Colour != red || cl.Len != 1 {
			return c, false
		}

		c.Len += cl.Len
		c.BlackH += cl.BlackH
		c.MinH += cl.MinH
		c.MaxH += cl.MaxH
	} else {
		cr, ok := n.right.check()
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

func (n *node) rangeAllAsc(res []pair.Pair, pos *int) {
	if n == nilNode {
		return
	}

	if n.left != nilNode {
		n.left.rangeAllAsc(res, pos)
	}
	res[*pos].First = n.key
	res[*pos].Second = n.value
	*pos++
	if n.right != nilNode {
		n.right.rangeAllAsc(res, pos)
	}
}

func (n *node) rangeKeysAsc(res []interface{}, pos *int) {
	if n == nilNode {
		return
	}
	if n.left != nilNode {
		n.left.rangeKeysAsc(res, pos)
	}
	res[*pos] = n.key
	*pos++
	if n.right != nilNode {
		n.right.rangeKeysAsc(res, pos)
	}
}

func (n *node) rangeValuesAsc(res []interface{}, pos *int) {
	if n == nilNode {
		return
	}
	if n.left != nilNode {
		n.left.rangeValuesAsc(res, pos)
	}
	res[*pos] = n.value
	*pos++
	if n.right != nilNode {
		n.right.rangeValuesAsc(res, pos)
	}
}

func (n *node) rangeAllDesc(res []pair.Pair, pos *int) {
	if n == nilNode {
		return
	}

	if n.right != nilNode {
		n.right.rangeAllDesc(res, pos)
	}
	res[*pos].First = n.key
	res[*pos].Second = n.value
	*pos++
	if n.left != nilNode {
		n.left.rangeAllDesc(res, pos)
	}
}

func (n *node) rangeAsc(res []pair.Pair, minKey, maxKey interface{}, cmp CmpFunc) []pair.Pair {
	if n == nilNode {
		return res
	}

	cmpMin, cmpMax := cmp(n.key, minKey), cmp(n.key, maxKey) // cmp() may takes some time, so we just cmp one time.
	if cmpMin > 0 {
		res = n.left.rangeAsc(res, minKey, maxKey, cmp)
	}
	if cmpMin >= 0 && cmpMax <= 0 {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if cmpMax < 0 {
		res = n.right.rangeAsc(res, minKey, maxKey, cmp)
	}
	return res
}

func (n *node) rangeAscN(res []pair.Pair, num int, key interface{}, cmp CmpFunc) []pair.Pair {
	if n == nilNode {
		return res
	}

	iCmp := cmp(n.key, key) // cmp() may takes some time, so we just cmp one time.
	if iCmp > 0 && len(res) < num {
		res = n.left.rangeAscN(res, num, key, cmp)
	}
	if iCmp >= 0 && len(res) < num {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if len(res) < num {
		res = n.right.rangeAscN(res, num, key, cmp)
	}
	return res
}

func (n *node) rangeDesc(res []pair.Pair, minKey, maxKey interface{}, cmp CmpFunc) []pair.Pair {
	if n == nilNode {
		return res
	}

	cmpMin, cmpMax := cmp(n.key, minKey), cmp(n.key, maxKey) // cmp() may takes some time, so we just cmp one time.
	if cmpMax < 0 {
		res = n.right.rangeDesc(res, minKey, maxKey, cmp)
	}
	if cmpMin >= 0 && cmpMax <= 0 {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if cmpMin > 0 {
		res = n.left.rangeDesc(res, minKey, maxKey, cmp)
	}
	return res
}

func (n *node) rangeDescN(res []pair.Pair, num int, key interface{}, cmp CmpFunc) []pair.Pair {
	if n == nilNode {
		return res
	}

	iCmp := cmp(n.key, key) // cmp() may takes some time, so we just cmp one time.
	if iCmp < 0 && len(res) < num {
		res = n.right.rangeDescN(res, num, key, cmp)
	}
	if iCmp <= 0 && len(res) < num {
		res = append(res, pair.Pair{First: n.key, Second: n.value})
	}
	if len(res) < num {
		res = n.left.rangeDescN(res, num, key, cmp)
	}
	return res
}

// O(N)
func (n *node) getDepth() uint {
	if n == nilNode {
		return 0
	}
	ld := n.left.getDepth()
	rd := n.right.getDepth()

	if ld > rd {
		return ld + 1
	} else {
		return rd + 1
	}
}

// O(logN)
func (n *node) getLeftDepth() uint {
	if n == nilNode {
		return 0
	}
	return n.left.getLeftDepth() + 1
}

// Deprecated: only for debugging, unstable function
func (n *node) makeString(buffs []*strings.Builder, step, lMove, tDepth, nDepth uint, ifRowFirst, ifParentRowFirst bool) {
	if n == nilNode {
		return
	}
	space := (step << (tDepth - nDepth)) - step
	if ifRowFirst {
		if space >= lMove {
			space -= lMove
		}
	} else {
		space = (space << 1) + step
		if ifParentRowFirst && n.parent.left == nilNode {
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

	if n.left != nilNode {
		n.left.makeString(buffs, step, lMove, tDepth, nDepth+1, ifRowFirst, false)
	} else if !ifRowFirst {
		for i := nDepth + 1; i <= tDepth; i++ {
			rStep := spaceCount(step, tDepth, i)
			rSpace := rStep << (i - nDepth - 1)
			for j := rSpace; j > 0; j-- {
				buffs[i-1].WriteString(" ")
			}
		}
	}

	if n.right != nilNode {
		if n.left == nilNode {
			n.right.makeString(buffs, step, lMove, tDepth, nDepth+1, false, true)
		} else {
			n.right.makeString(buffs, step, lMove, tDepth, nDepth+1, false, false)
		}
	} else {
		if ifRowFirst && n.left == nilNode {
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
