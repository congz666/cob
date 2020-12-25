package cob

import "strings"

type node struct {
	pattern  string  // Routes to be matched. This value is zero for all nodes except the lowest node. e.g. /p/:lang
	part     string  // Part of a route, e.g. :lang
	children []*node // Child nodes，e.g. [doc, tutorial, intro]
	isWild   bool    // Whether to match dynamically. True when part contains : or *
}

// For insert
// Find the first node that matches successfully
func (n *node) matchChildByInsert(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// For search
// Find the all nodes that matches successfully
func (n *node) matchChildByFind(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// Detect conflicts
func (n *node) matchChildByJudge(part string) ([]*node, bool) {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.isWild == true {
			return nil, false
		}
		if child.part == part {
			nodes = append(nodes, child)
		}
	}
	return nodes, true
}

// Insert route
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChildByInsert(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// Match route
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	part := parts[height]
	children := n.matchChildByFind(part)

	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}

	return nil
}

// Determine whether there is a routing conflict
func (n *node) judge(parts []string, height int) (*node, bool) {
	if len(parts) == height {
		return nil, false
	}

	part := parts[height]
	children, ok := n.matchChildByJudge(part)
	if !ok {
		return nil, false
	}

	for _, child := range children {
		result, ok := child.judge(parts, height+1)
		if !ok {
			return nil, false
		}
		if result != nil {
			return result, true
		}
	}

	return nil, true
}
