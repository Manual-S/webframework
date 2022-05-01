package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	isLast  bool
	segment string
	handler ControllerHandler
	childs  []*node
}

func newNode() *node {
	return nil
}

// isWildSegment 判断一个segment是否是通配segment 以:开头的就是通配segment
func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

// filterChildNodes 过滤下一层满足segment规则的子节点
func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}

	if isWildSegment(segment) {
		// 说明segment是通配符 所有下层节点都要满足
		return n.childs
	}

	nodes := make([]*node, 0)
	for _, cnode := range n.childs {
		if isWildSegment(cnode.segment) {
			// 如果下一层子节点含有通配符 则满足要求
			nodes = append(nodes, cnode)
		} else if cnode.segment == segment {
			nodes = append(nodes, cnode)
		}
	}

	return nodes
}

func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]

	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}

	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		// 说明uri segment已经是最后一个节点
	}

	// 如果有2个segment 递归每个子节点继续进行查找
	for _, tn := range cnodes {
		tnMatch := tn.matchNode(segments[1])
		if tnMatch != nil {
			return tnMatch
		}
	}

	return nil
}

// AddRouter 增加路由规则
func (tree *Tree) AddRouter(uri string, handler ControllerHandler) error {
	n := tree.root
	if n.matchNode(uri) != nil {
		return errors.New("roter exist" + uri)
	}

	segments := strings.Split(uri, "/")

	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}
		var isLast bool
		if index == len(segments)-1 {
			isLast = true
		} else {
			isLast = false
		}

		var objNode *node

		childNodes := n.filterChildNodes(segment)

		if len(childNodes) > 0 {
			// 说明有匹配的子节点
		}

		if objNode == nil {
			// 说明没有匹配的子节点 创建一个新的子节点 挂载
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}

			n.childs = append(n.childs, cnode)

			objNode = cnode
		}

		n = objNode
	}

	return nil
}
