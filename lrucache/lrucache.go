package lrucache

import (
	"fmt"
	"sync"
)

type Node struct {
	prev *Node
	next *Node
	val  int
	key  int
}

type LRU struct {
	right *Node
	left  *Node
	cap   int
	cache map[int]*Node
	my    sync.RWMutex
}

func NewLruCLient(cap int) *LRU {
	lru := &LRU{
		right: &Node{},
		left:  &Node{},
		cap:   cap,
		cache: make(map[int]*Node),
	}
	lru.right.prev = lru.left
	lru.left.next = lru.right
	return lru
}

func (lru *LRU) Put(key, value int) bool {
	lru.my.Lock()
	defer lru.my.Unlock()

	if n, exists := lru.cache[key]; exists {
		fmt.Printf("Already exists, updating %v with %v \n", key, value)
		lru.moveNodeToHead(n)
		n.val = value
		return true
	}

	fmt.Printf("Creating new key %v with %v \n", key, value)
	newNode := &Node{key: key, val: value}
	lru.addNode(newNode)
	lru.cache[key] = newNode

	if len(lru.cache) > lru.cap {
		lruNode := lru.left.next
		fmt.Printf("Deleting LRU as max cap reached %v \n", lruNode.key)
		delete(lru.cache, lruNode.key)
		removeNode(lruNode)
	}

	return true
}

func (lru *LRU) Get(key int) int {
	lru.my.RLock()
	node, exists := lru.cache[key]
	lru.my.RUnlock()
	if !exists {
		return -1
	}

	// Now we need to reorder -> exclusive lock
	lru.my.Lock()
	lru.moveNodeToHead(node)
	lru.my.Unlock()

	return node.val
}

func removeNode(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// head is left, tail is right
func (lru *LRU) addNode(n *Node) {
	currRight := lru.right
	currPrev := lru.right.prev

	currPrev.next = n
	currRight.prev = n

	n.prev = currPrev
	n.next = currRight
}

func (lru *LRU) moveNodeToHead(node *Node) {
	removeNode(node)
	lru.addNode(node)
}

func Run() {
	LRUClient := NewLruCLient(2)

	res := LRUClient.Get(1)
	fmt.Printf("Get KEY: 1 VALUE %v \n", res)

	LRUClient.Put(1, 100)
	LRUClient.Put(1, 6969)

	res = LRUClient.Get(1)
	fmt.Printf("Get KEY: 1 VALUE %v\n", res)

	LRUClient.Put(3, 300)
	LRUClient.Put(4, 400)

	res = LRUClient.Get(1)
	fmt.Printf("Get KEY: 1 VALUE %v\n", res)

	res = LRUClient.Get(3)
	fmt.Printf("Get KEY: 3 VALUE %v\n", res)

}
