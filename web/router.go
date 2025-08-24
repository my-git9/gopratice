package web

import "strings"

// 用来支持对路由树的操作
// 代表路由树（森林）
type router struct {
	// Beego Gin HTTP method 对应一棵树
	// Get 一棵树， Post 一棵树

	// http methods ==> 路由树根节点
	trees map[string]*node
}

func newRouter() *router {
	return &router{
		trees: map[string]*node{

		},
	}
}

// TODO: path 必须以 / 开头，不能以 / 结尾，中间也不能有连续的 //
// TODO: 重复注册
func (r *router)  AddRoute(method string, path string, handleFunc HandleFunc)  {
	if path == "" {
		panic("web： 路径不能为空字符事")
	}
	// 首先找到树
	root, ok := r.trees[method]
	if !ok {
		// 说明还没有根节点
		root := &node{
			path: "/",
		}
		r.trees[method] = root
	}

	// 开头不能没有 /
	if path[0] != '/' {
		panic("web： 开头不能没有 /")
	}

	// 根节点特殊处理
	if path == "/" {
		root.handler = handleFunc
		return
	}

	// 切割这个 path
	path = path[1:] // 去掉开头的 /
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		// 不能有 //
		if seg == "" {
			panic("不能有 //")
		}
		// 递归下去，找准位置
		// 如果中途有节点不存在，就要创建出来
		children := root.ChildOrCreate(seg)
		root = children
	}
}

type node struct {
	path string

	// 子 path 到 子节点的映射
	children map[string]*node

	// 缺一个代表用户注册的业务逻辑
	handler HandleFunc
}

func (r *router) findRoute(method string, path string) (*node, bool) {
	root, ok := r.trees[method]
	if !ok {
		return nil, false
	}
	// 这里把前置和后置的 / 去掉
	path = strings.Trim(path, "/")
	// 按照 / 切割
	segs := strings.Split(path, "/")
	for _, seg := range segs {
		child, found := root.ChildOf(seg)
		if !found {
			return nil, false
		}
		root = child
	}
	// 代表确实有这个节点
	// 但是这个节点是不是注册有 handler 的，就不一定（留给使用者自行判断）
	return root, true
	//return root, root.handler != nil
}

// 第一个返回值是正确的子节点
// 第二个
func (n *node) ChildOrCreate(seg string) *node {
	if n.children == nil {
		res := &node{
			path: seg,
		}
		n.children[seg] = res
		return res
	}
	res, ok := n.children[seg]
	if !ok {
		// 新建一个
		res = & node {
			path: seg,
		}
		n.children[seg] = res
	}
	return res
}

func (n *node) ChildOf(path string) (*node, bool) {
	if n.children == nil {
		return nil, false
	}
	child, ok := n.children[path]
	return child, ok
}

