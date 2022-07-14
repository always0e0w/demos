package main

type Node struct {
	pre, next *Node
	key, val  string
}

func newNode(k, v string) *Node {
	return &Node{
		key: k,
		val: v,
	}
}

type LRU struct {
	capacity, size int
	head, tail     *Node // dummy node
	keyMap         map[string]*Node
}

func NewLRU(n int) *LRU {
	lru := &LRU{
		capacity: n,
		keyMap:   make(map[string]*Node),
		head:     newNode("", ""),
		tail:     newNode("", ""),
	}
	lru.head.pre = lru.tail
	lru.head.next = lru.tail
	lru.tail.pre = lru.head
	lru.tail.next = lru.head
	return lru
}

func (lru LRU) Get(k string) (string, bool) {
	node, e := lru.keyMap[k]
	if !e {
		return "", false
	}
	lru.moveToHead(node)
	return node.val, true
}

func (lru *LRU) moveToHead(node *Node) {
	lru.removeNode(node)
	lru.addNode(node)
}

func (lru *LRU) Set(k, v string) {
	if node, e := lru.keyMap[k]; e {
		node.val = v
		lru.moveToHead(node)
		return
	}
	node := newNode(k, v)
	lru.addNode(node)
	lru.keyMap[k] = node
	lru.size += 1
	if lru.size > lru.capacity {
		node = lru.removeLast()
		delete(lru.keyMap, node.key)
		lru.size -= 1
	}
}

func (lru *LRU) removeLast() *Node {
	if lru.size == 0 {
		return nil
	}
	node := lru.tail.pre
	lru.removeNode(node)
	return node
}

func (lru *LRU) removeNode(node *Node) {
	if node == nil {
		return
	}
	node.pre.next = node.next
	node.next.pre = node.pre
	node.pre = nil
	node.next = nil
}

func (lru *LRU) addNode(node *Node) {
	if node == nil {
		return
	}
	node.pre = lru.head
	node.next = lru.head.next
	lru.head.next.pre = node
	lru.head.next = node
}
