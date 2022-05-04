package framework

import (
	"errors"
	"fmt"
	"strings"
)

type Tree struct {
	root *node
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{
		root: root,
	}
}

type node struct {
	isLast  bool                // 代表这个节点是否可以成为最终的路由规则
	segment string              // uri中的字符串
	handler []ControllerHandler // 中间件+控制器
	childs  []*node             // 子节点
	parent  *node               // 父节点
}

func newNode() *node {
	return &node{
		isLast:  false,
		segment: "",
		childs:  []*node{},
	}
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

// matchNode 找到当前路由匹配的节点
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]

	if !isWildSegment(segment) {
		// 把所有的uri改为大写
		segment = strings.ToUpper(segment)
	}

	// 匹配符合的下一层子节点
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	if len(segments) == 1 {
		// 说明uri segment已经是最后一个节点
		for _, tn := range cnodes {
			if tn.isLast {
				return tn
			}
		}
		return nil
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
func (tree *Tree) AddRouter(uri string, handler ...ControllerHandler) error {
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
			for _, cnode := range childNodes {
				if cnode.segment == segment {
					objNode = cnode
					break
				}
			}
		}

		if objNode == nil {
			// 说明没有匹配的子节点 创建一个新的子节点 挂载
			cnode := newNode()
			cnode.segment = segment
			if isLast {
				cnode.isLast = true
				cnode.handler = handler
			}

			cnode.parent = n

			n.childs = append(n.childs, cnode)

			objNode = cnode
		}

		n = objNode
	}

	return nil
}

// FindHandler 匹配uri
func (tree *Tree) FindHandler(uri string) *node {
	matchNode := tree.root.matchNode(uri)
	if matchNode == nil {
		return nil
	}
	return matchNode
}

// 将uri解析成params
func (n *node) parseParmsFromEndNode(uri string) map[string]string {
	ret := map[string]string{}

	segment := strings.Split(uri, "/")

	cnt := len(segment)
	cur := n

	for i := cnt - 1; i >= 0; i-- {
		fmt.Printf("cur.segment = %v", cur.segment)
		if cur.segment == "" {
			break
		}

		if isWildSegment(cur.segment) {
			// 说明是一个通配符
			ret[cur.segment[1:]] = segment[i]
		}
		cur = cur.parent
	}

	return ret
}
