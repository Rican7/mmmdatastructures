// Package stringqueue implements a queue for strings.
//
// The internal representation is a slice of strings
// that gets used as a circular buffer.
// This is instead of a more traditional approach
// that would use a linked list of nodes.
// The assumption is that contiguous slabs of RAM
// will generally provide more performance over pointers
// to nodes potentially scattered about the heap.
//
// There is a downside: whereas enqueueing to a
// linked list is always O(1), enqueueing here will
// be O(1) except for when the internal slice of strings
// has to be resized; then, enqueueing will be O(n)
// where n is the size of the queue before being resized.
//
// Therefore, when asking for a new instance of the
// queue, pick a capacity that you think won't need to grow.
//
// When the queue does need to grow, it always uses a capacity
// that is twice the current capacity. This is not tunable.
package stringqueue

import "github.com/pkg/errors"

// IntQueue holds the data and state of the queue.
type IntQueue struct {
	data     []string
	head     int
	tail     int
	capacity int
	length   int
}

// DefauiltCapacity is the default capacity of the IntQueue
// when constructed using New() instead of NewWithCapacity().
const DefaultCapacity = 32

// New returns a new empty queue for strings of the default capacity.
func New() (q *IntQueue) {
	return NewWithCapacity(DefaultCapacity)
}

// NewWithCapacity returns a new empty queue for strings with the requested capacity.
func NewWithCapacity(capacity int) (q *IntQueue) {
	q = new(IntQueue)
	q.data = make([]string, capacity, capacity)
	q.head = -1
	q.tail = -1
	q.capacity = capacity
	q.length = 0
	return q
}

// Enqueue enqueues a string. Returns an error if the size
// of the queue cannot be grown any more to accommodate
// the added string.
func (q *IntQueue) Enqueue(i string) error {
	if q.length+1 > q.capacity {
		new_capacity := q.capacity * 2
		// if new_cap became negative, we have exceeded
		// our capacity by doing one bit-shift too far
		if new_capacity < 0 {
			return errors.New("Capacity exceeded")
		}
		q.resize(new_capacity)
	}
	q.length++
	q.head++
	if q.head == q.capacity {
		q.head = 0
	}
	q.data[q.head] = i
	return nil
}

// Head can be earlier in array than tail, so
// we can't just copy; we could overwrite the tail.
// Instead, we may as well copy the queue in order
// into the new array. The Dequeue() method gives us
// every element in the correct order already, so we
// just leverage that.
func (q *IntQueue) resize(new_capacity int) {
	new_data := make([]string, new_capacity, new_capacity)
	var err error
	var i string
	for err = nil; err == nil; i, err = q.Dequeue() {
		new_data = append(new_data, i)
	}
	q.head = q.length - 1
	q.tail = 0
	q.capacity = new_capacity
	q.data = new_data
}

// Dequeue dequeues a string. It returns the dequeued string
// or an error of the queue is empty.
func (q *IntQueue) Dequeue() (string, error) {
	if q.length-1 < 0 {
		return "", errors.New("Queue empty")
	}
	q.length--
	q.tail++
	if q.tail == q.capacity {
		q.tail = 0
	}
	return q.data[q.tail], nil
}