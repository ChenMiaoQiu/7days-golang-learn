package gee

import "strings"

type node struct {
	path     string
	part     string
	children []*node
	isWild   bool
}

func (n *node) mathChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

func (n *node) mathChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(path string, parts []string, height int) {
	if len(parts) == height {
		n.path = path
		return
	}

	part := parts[height]
	child := n.mathChild(part)
	if child == nil {
		child = &node{part: part, isWild: (part[0] == ':' || part[0] == '*')}
		n.children = append(n.children, child)
	}
	child.insert(path, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.path == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.mathChildren(part)

	for _, child := range children {
		res := child.search(parts, height+1)
		if res != nil {
			return res
		}
	}

	return nil
}
