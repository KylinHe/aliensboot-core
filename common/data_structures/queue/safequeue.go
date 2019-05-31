package queue

import "sync"

type SafeQueue struct {
	queue []interface{}
	len   int
	lock  *sync.Mutex
}

func NewSafe() *SafeQueue {
	queue := &SafeQueue{}
	queue.queue = make([]interface{}, 0)
	queue.len = 0
	queue.lock = new(sync.Mutex)

	return queue
}

func (q *SafeQueue) Len() int {
	//q.lock.Lock()
	//defer q.lock.Unlock()

	return q.len
}

func (q *SafeQueue) isEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.len == 0
}

func (q *SafeQueue) Pop() (el interface{}) {
	//q.lock.Lock()
	//defer q.lock.Unlock()

	el, q.queue = q.queue[0], q.queue[1:]
	q.len--
	return
}

type RuleFunc func(interface{}) bool

func (q *SafeQueue) PopRule(rule RuleFunc) (el interface{}) {
	//q.lock.Lock()
	//defer q.lock.Unlock()
	for index := 0; index < q.len; index++ {
		if rule(q.queue[index]) {
			el = q.queue[index]
			q.queue = append(q.queue[:index], q.queue[index+1:]...)
			q.len--
			break
		}
	}

	return
}

func (q *SafeQueue) Push(el interface{}) {
	//q.lock.Lock()
	//defer q.lock.Unlock()

	q.queue = append(q.queue, el)
	q.len++

	return
}

func (q *SafeQueue) Peek() interface{} {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.queue[0]
}
