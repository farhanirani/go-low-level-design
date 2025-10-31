// LLD Walkthrough [Golang]: LRU Cache

// ----------------------------------------------------------------------
// Design an in-memory LRU (Least Recently Used) Cache with fixed capacity.
// The cache should evict the least recently used key when the capacity is reached.
// It should support concurrent access from multiple goroutines.

// ******** STEP 1. REQUIREMENTS
// 1. get(key): return the value if key exists; else return -1.
// 2. put(key, value): insert or update a key-value pair.
//    - If the cache exceeds capacity, remove the least recently used key.
// 3. The cache must maintain usage order (most recently used near tail).
// 4. The cache must be thread-safe (support concurrent access).
// 5. Both get and put operations should run in O(1) time complexity.

// Q: How do we achieve O(1) get and put operations?
// A: Combine a hash map (for fast key lookup) and a doubly linked list (for ordering).

// ----------------------------------------------------------------------

package lrucache

import (
	"fmt"
	"sync"
)

// ******** STEP 2 - DISCUSS ENTITIES

// Node → represents one cache entry (key-value pair) in the doubly linked list
// Each node also has pointers to its previous and next nodes
type Node struct {
	prev *Node
	next *Node
	val  int
	key  int
}

// LRU → represents the entire LRU Cache
// Maintains:
// - doubly linked list (for usage order)
// - hash map for O(1) access
// - capacity
// - RWMutex for concurrency control
type LRU struct {
	right *Node         // dummy head (most recently used)
	left  *Node         // dummy tail (least recently used)
	cap   int           // cache capacity
	cache map[int]*Node // key → node pointer
	my    sync.RWMutex  // protects cache and linked list
}

// ******** STEP 3 - FACTORY & CORE METHODS / APIs

// put(key, value) → Insert or update a key-value pair
// get(key) → Retrieve value; move it to most recently used
// Thread-safe via read/write locks

// ----------------------------
// FACTORY FUNCTION
// ✅ Factory Method Pattern — encapsulates initialization logic
// ----------------------------
func NewLruCLient(cap int) *LRU {
	lru := &LRU{
		right: &Node{}, // dummy head
		left:  &Node{}, // dummy tail
		cap:   cap,
		cache: make(map[int]*Node),
	}

	// Initialize pointers
	lru.right.prev = lru.left
	lru.left.next = lru.right

	return lru
}

// ----------------------------
// CORE METHODS
// ----------------------------

// Put inserts or updates a key-value pair in O(1)
func (lru *LRU) Put(key, value int) bool {
	lru.my.Lock()
	defer lru.my.Unlock()

	// Case 1: key already exists — update and move to head
	if n, exists := lru.cache[key]; exists {
		fmt.Printf("Already exists, updating %v with %v \n", key, value)
		n.val = value
		lru.moveNodeToHead(n)
		return true
	}

	// Case 2: new key — create new node
	fmt.Printf("Creating new key %v with %v \n", key, value)
	newNode := &Node{key: key, val: value}
	lru.addNode(newNode)
	lru.cache[key] = newNode

	// Case 3: capacity overflow — evict least recently used
	if len(lru.cache) > lru.cap {
		lruNode := lru.left.next
		fmt.Printf("Deleting LRU as max cap reached %v \n", lruNode.key)
		delete(lru.cache, lruNode.key)
		removeNode(lruNode)
	}

	return true
}

// Get retrieves value in O(1)
// Moves accessed node to most recently used (tail)
func (lru *LRU) Get(key int) int {
	lru.my.RLock()
	node, exists := lru.cache[key]
	lru.my.RUnlock()

	if !exists {
		return -1
	}

	// Need exclusive lock to reorder the node
	lru.my.Lock()
	lru.moveNodeToHead(node)
	lru.my.Unlock()

	return node.val
}

// ----------------------------
// INTERNAL LINKED LIST OPERATIONS
// ----------------------------

// removeNode removes a node from the doubly linked list in O(1)
func removeNode(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// addNode adds a node to the tail (most recently used position)
func (lru *LRU) addNode(n *Node) {
	currRight := lru.right
	currPrev := lru.right.prev

	currPrev.next = n
	currRight.prev = n

	n.prev = currPrev
	n.next = currRight
}

// moveNodeToHead moves an existing node to the most recently used position
func (lru *LRU) moveNodeToHead(node *Node) {
	removeNode(node)
	lru.addNode(node)
}

// ----------------------------
// DEMO
// ----------------------------
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

// ----------------------------------------------------------------------
// ******** STEP 4 - ANALYSIS / EXPLANATION
// ----------------------------------------------------------------------
// ✅ Data Structures Used:
// - map[int]*Node for O(1) key lookups
// - doubly linked list for O(1) insertion/deletion and maintaining order

// ✅ Order Convention:
// - left.next = least recently used
// - right.prev = most recently used

// ✅ Concurrency Control:
// - RWMutex allows multiple readers, single writer
// - Put() uses Lock() since it modifies cache
// - Get() uses RLock() for lookup, then Lock() for node movement

// ✅ Complexity:
// - Time: O(1) for both Get and Put
// - Space: O(capacity) for nodes + map

// ----------------------------------------------------------------------
// ******** STEP 5 - POSSIBLE FOLLOW-UP QUESTIONS
// ----------------------------------------------------------------------

/*
Question / Follow-up                        Concept / Pattern              High-Level Answer

How to make this highly concurrent?         Sharding / Partitioned Locking Divide the cache into multiple shards (e.g., by hash(key) % N),
                                                                           each with its own lock, to reduce contention.

How to persist cache state on restart?      Serialization / Snapshotting   Store the map to disk using Gob or JSON encoding and restore on startup.

How to add TTL expiry?                      Background Goroutine / Heap    Maintain a min-heap of expiry times; periodically evict expired entries.

Can this cache be distributed?              Consistent Hashing             Use consistent hashing to distribute keys across nodes (like memcached).

How to reduce GC overhead?                  sync.Pool / Object Reuse       Use sync.Pool for Node allocation to reduce memory churn.

Why use doubly linked list over slice?      Time complexity tradeoff       Slice deletion in the middle is O(n); doubly linked list makes it O(1).

Can reads be fully lock-free?               Atomic Pointer Swapping        Use atomic.Value or copy-on-write for read-mostly workloads.
*/

// *************************************************
// REMEMBER: In interviews, correctness + clarity >>> fancy optimizations
// Explain trade-offs before you code optimizations.
// *************************************************

// NOTES:
// Concept                     Equivalent / Analogy
// map[int]*Node               Dictionary lookup for O(1) access
// Doubly linked list          Keeps track of recency order efficiently
// sync.RWMutex                Reader-writer lock for thread safety
// removeNode + addNode        Core LRU operations maintaining order
