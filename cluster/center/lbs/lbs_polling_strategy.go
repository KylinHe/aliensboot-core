/*******************************************************************************
 * Copyright (c) 2015, 2018 aliens idea(xiamen) Corporation and others.
 * All rights reserved.
 * Date:
 *     2017/11/3
 * Contributors:
 *     aliens idea(xiamen) Corporation - initial API and implementation
 *     jialin.he <kylinh@gmail.com>
 * Desc:
 *     Load Balance Strategy -- polling
 *******************************************************************************/
package lbs

func NewPollingLBS() *PollingLBS {
	return &PollingLBS{nodes: make(map[string]struct{})}
}

//轮询负载策略
type PollingLBS struct {
	nodes    map[string]struct{}
	rootNode *PollingNode //链表根节点
	currNode *PollingNode //当前请求到的节点
}

type PollingNode struct {
	name   string
	weight int
	next   *PollingNode
	index  int
}

//func (this *PollingLBS) Init(nodes []string) {
//	this.nodes = nodes
//	this.length = len(this.nodes)
//}

func (this *PollingLBS) AddNode(nodeKey string, weight int) {
	_, ok := this.nodes[nodeKey]
	if ok {
		return
	}
	this.nodes[nodeKey] = struct{}{}

	newRootNode := &PollingNode{name: nodeKey, weight: weight, index: 0}
	if this.rootNode != nil {
		newRootNode.next = this.rootNode
	} else {
		this.currNode = newRootNode
	}
	this.rootNode = newRootNode
}

func (this *PollingLBS) RemoveNode(nodeKey string) {
	_, ok := this.nodes[nodeKey]
	if !ok {
		return
	}
	delete(this.nodes, nodeKey)

	node := this.rootNode
	if node == nil {
		return
	}
	var lastNode *PollingNode = nil

	for {
		if node == nil {
			return
		}
		if node.name == nodeKey {
			//当前节点为删除节点、需要切换到下一个节点
			if lastNode != nil {
				lastNode.next = node.next
			}
			if this.rootNode == node {
				this.rootNode = node.next
			}
			if this.currNode == node {
				this.nextNode()
			}
			return
		} else {
			lastNode = node
			node = node.next
		}
	}
}

func (this *PollingLBS) GetNode(key string) string {
	if this.currNode == nil {
		return ""
	}
	this.currNode.index++
	//超过权重 取下个节点
	if this.currNode.index > this.currNode.weight {
		this.nextNode()
	}
	return this.currNode.name

}

func (this *PollingLBS) nextNode() {
	this.currNode.index = 0
	this.currNode = this.currNode.next
	//链表到了最后一个节点，重新取第一个节点
	if this.currNode == nil {
		this.currNode = this.rootNode
	}
}
