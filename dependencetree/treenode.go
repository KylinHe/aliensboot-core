package dependencetree

import (
	"github.com/KylinHe/aliensboot-core/log"
)

type Tree struct {

	indies map[string]*TreeNode  //所有节点索引

}

//func (t *Tree) RegisterDepends(parent ITreeNode, child ITreeNode) {
//	node := t.indies[parent.GetID()]
//	if node == nil {
//		node = &TreeNode{
//			ITreeNode:parent,
//		}
//	}
//}

func (t *Tree) SetCustomData(category Category, id int32, customData interface{}) {
	node := t.EnsureNode(category, id)
	node.customData = customData
}

// 注册依赖
func (t *Tree) RegisterDepend(category Category, id int32, dependCategory Category, dependId int32) {
	node := t.EnsureNode(category, id)
	dependNode := t.EnsureNode(dependCategory, dependId)
	node.addDepend(dependNode)
	dependNode.addTrigger(node)
}

// 按层级获取依赖节点
func (t *Tree) GetDependsByLayer(category Category, id int32, layer int, filter func(node *TreeNode, layer int) bool) map[string]*TreeNode {
	//todo 按层级获取依赖节点
	node := t.EnsureNode(category, id)
	result := make(map[string]*TreeNode)
	node.GetDepends(result, layer, 1, filter)
	return result
}

// 按层级获取触发节点
func (t *Tree) GetTriggersByLayer(category Category, id int32, layer int, filter func(node *TreeNode, layer int) bool) map[string]*TreeNode {
	//todo 按层级获取触发节点
	node := t.EnsureNode(category, id)
	result := make(map[string]*TreeNode)
	node.GetTriggers(result, layer, 1, filter)
	return result
}

func (t *Tree) EnsureNode(category Category, id int32) *TreeNode {
	cateId := CategoryId(category, id)
	node := t.indies[cateId]
	if node == nil {
		node = NewTreeNode(category, id)
		t.indies[cateId] = node
	}
	return node
}

func CategoryId(category Category, id int32) string {
	return string(category) + string(id)
}

type Category string

func NewTreeNode(category Category, id int32) *TreeNode {
	return &TreeNode{
		categoryId:CategoryId(category,id),
		id:id,
		category:category,
	}
}

type TreeNode struct {

	categoryId string

	category Category

	id int32

	depends map[string]*TreeNode   //依赖的节点

	triggers map[string]*TreeNode  //触发的节点

	customData interface{}
}

func (node *TreeNode) GetCustomData() interface{} {
	return node.customData
}

func (node *TreeNode) GetDepends(results map[string]*TreeNode, maxLayer int, layer int, filter func(node *TreeNode, layer int) bool) {
	if layer > maxLayer {
		return
	}
	for key, node := range node.depends {
		if filter(node, layer) {
			continue
		}
		results[key] = node
	}

	layer++
	for _, node := range node.depends {
		node.GetDepends(results, maxLayer, layer, filter)
	}
}

func (node *TreeNode) GetTriggers(results map[string]*TreeNode, maxLayer int, layer int, filter func(node *TreeNode, layer int) bool) {
	if layer > maxLayer {
		return
	}
	for key, node := range node.triggers {
		if filter(node, layer) {
			continue
		}
		results[key] = node
	}

	layer++
	for _, node := range node.triggers {
		node.GetDepends(results, maxLayer, layer, filter)
	}
}

func (node *TreeNode) addDepend(dependNode *TreeNode) {
	if node.depends == nil {
		node.depends = make(map[string]*TreeNode)
	}
	oldDependNode := node.depends[dependNode.categoryId]
	if oldDependNode != nil {
		log.Debugf("depend node already exist %v => %v", dependNode.categoryId, node.categoryId)
	} else {
		node.depends[dependNode.categoryId] = dependNode
	}
}

func (node *TreeNode) addTrigger(triggerNode *TreeNode) {
	if node.triggers == nil {
		node.triggers = make(map[string]*TreeNode)
	}
	oldTriggerNode := node.triggers[triggerNode.categoryId]
	if oldTriggerNode != nil {
		log.Debugf("trigger node already exist %v = %v", node.categoryId, triggerNode.categoryId)
	} else {
		node.triggers[triggerNode.categoryId] = triggerNode
	}
}