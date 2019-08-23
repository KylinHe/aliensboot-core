package dependencetree


type Tree struct {
	indies map[string]*TreeNode  //节点索引
}

func (t *Tree) RegisterDepends(parent ITreeNode, child ITreeNode) {
	node := t.indies[parent.GetID()]
	if node == nil {
		node = &TreeNode{
			ITreeNode:parent,
		}
	}
}


type ITreeNode interface {
	GetID() string
	IsOpen() bool
	IsDone()
}

type TreeNode struct {
	ITreeNode

	parents map[string]*TreeNode   //
	children map[string]*TreeNode  //
}
