package gee

import (
	"strings"
)

type node struct {
	pattern string
	part string
	children []*node
	isWild bool
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	part := parts[height]
	for _, child := range n.children {
		if child.part == part || child.isWild {
			child.insert(pattern, parts, height + 1)
			return
		}
	}
	child := &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
	n.children = append(n.children, child)
	child.insert(pattern, parts, height + 1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	for _, child := range n.children {
		if child.part == part || child.isWild {
			result := child.search(parts, height + 1)
			if result != nil {
				return result
			}
		}
	}
	return nil
}

func (n *node) travel(list *([]*node)) {
	if n.pattern != "" {
		*list = append(*list, n)
	}
	for _, child := range n.children {
		child.travel(list)
	}
}