package main

import "crypto/sha256"

type MerkleTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left *MerkleNode
	Rright *MerkleNode
	Data []byte
}

func NewMerkleNode(left *MerkleNode,right *MerkleNode, data []byte) *MerkleNode{
	mnode := MerkleNode{}
	if left == nil && right== nil {
		mnode.Data = data
	}else {
		prehashes := append(left.Data,right.Data...)
		firsthash := sha256.Sum256(prehashes)
		hash:=sha256.Sum256(firsthash[:])
		mnode.Data=hash[:]
	}
	mnode.Left = left
	mnode.Rright = right
	return &mnode
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode
	//叶子节点
	for _,datum := range data{
		node := NewMerkleNode(nil, nil,datum)
		nodes = append(nodes, *node)
	}

	j:=0
	for nSize := len(data);nSize>1;nSize = (nSize + 1)/2 {
		for i:=0; i < nSize; i+=2 {
			//i2  为了当个数为奇数的时候，拷贝最后的元素
			i2 := min(i+ 1,nSize - 1)
			node := NewMerkleNode(&nodes[j+i], &nodes[j+i2], nil)
			nodes = append(nodes, *node)
		}
		j += nSize
	}
	mTree := MerkleTree{ & (nodes[len(nodes)-1])}
	return &mTree
}

