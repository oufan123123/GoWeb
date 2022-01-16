package gee

import "strings"

type node struct {
	pattern  string  // route pattern
	part     string  // route pattern part
	children []*node // child
	isWild   bool    // part contain : * isWild = true
}

var root node

// match one child used for insert
func (n *node) matchOne(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// match all childen used for search
func (n *node) matchAll(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	child := n.matchOne(parts[height])
	// create node
	if child == nil {
		wild := false
		if parts[height][0] == ':' || parts[height][0] == '*' {
			wild = true
		}
		child = &node{
			part:     parts[height],
			children: make([]*node, 0),
			isWild:   wild,
		}
		if height == len(parts) {
			child.pattern = pattern
		}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	// found
	if height == len(parts) || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]

	// all child
	nodes := n.matchAll(part)

	// not found
	if len(nodes) == 0 {
		return nil
	}
	for _, child := range nodes {
		found := child.search(parts, height+1)
		if found != nil {
			return found
		}
	}

	return nil
}
