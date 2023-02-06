package web

import "strings"

type node struct {
	pattern  string  //当插入结束时有值，标识着匹配路径，类似结束标志位
	path     string  //当前节点匹配的部分路径信息
	isWild   bool    // 标识当前节点是否进行模糊匹配
	children []*node //tries树中的子节点
}

func (n *node) matchChild(path string) *node {
	for _, v := range n.children {
		if v.path == path || v.isWild {
			return v
		}
	}
	return nil
}

// 和matchChild不同，对于查找匹配路径时在当前节点可能存在多个，此时需要对每个匹配项进行下一次筛选
func (n *node) matchChildren(path string) []*node {
	nodes := make([]*node, 0)
	for _, v := range n.children {
		if v.path == path || v.isWild {
			nodes = append(nodes, v)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, paths []string, height int) {
	if height == len(paths) {
		n.pattern = pattern
		return
	}
	path := paths[height]
	child := n.matchChild(path)
	if child == nil {
		child = &node{path: path, isWild: path[0] == ':' || path[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, paths, height+1)
}

func (n *node) search(paths []string, height int) *node {
	if height == len(paths) || strings.HasPrefix(n.path, "*") {
		// 如果pattern为"", 代表不能完全匹配，比如使用/a/b去匹配/a/b/c
		if n.pattern == "" {
			return nil
		}
		return n
	}
	path := paths[height]
	children := n.matchChildren(path)
	for _, v := range children {
		result := v.search(paths, height + 1)
		if result != nil {
			return result
		}
	}
	return nil
}
