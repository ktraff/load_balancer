package lib

import (
	"container/heap"
	"fmt"
	"github.com/ktraff/load_balancer"
)

// A list of workers.  We'll use a heap to respond to requests based
// on the worker that is the least busy.
type Pool []*Worker

type Balancer struct {
	pool Pool
	// Used to update the worker when a request has been processed
	done chan *Worker
}

func NewBalancer(numWorkers int, requestBufferSize int) *Balancer {
	fmt.Println(fmt.Sprintf("Creating balancers with %v workers (%v requests per worker)", numWorkers, requestBufferSize))
	done := make(chan *Worker, numWorkers)
	balancer := &Balancer{
		pool: make(Pool, 0, numWorkers),
		done: done,
	}
	backends := load_balancer.GetBackends()
	if len(*backends) == 0 {
		panic("No backends have been configured!")
	}
	for i := 1; i <= numWorkers; i++ {
		// Randomly assign a backend to a worker.
		backend := (*backends)[(i-1)%len(*backends)]
		worker := NewWorker(i, requestBufferSize, backend)
		fmt.Println(fmt.Sprintf("Creating %v", worker))
		heap.Push(&balancer.pool, &worker)
		go worker.work(balancer.done)
	}
	return balancer
}

func (b *Balancer) Balance(work chan *Request) {
	fmt.Println("Balancing load requests")
	for {
		select {
		case req := <-work:
			fmt.Println(fmt.Sprintf("Dispatching request: %v", req))
			b.dispatch(req)
		case w := <-b.done:
			fmt.Println(fmt.Sprintf("Finished request on %v", w))
			b.completed(w)
		}
	}
}

func (p *Pool) Len() int {
	return len(*p)
}

func (p *Pool) Push(x interface{}) {
	// Add an element to the end of the array
	pool := (*p)[0 : len(*p)+1]
	worker := x.(*Worker)
	pool[len(pool)-1] = worker
	worker.index = len(pool) - 1
	*p = pool
}

func (p *Pool) Pop() interface{} {
	// Remove an element from the end of the array
	last_worker := (*p)[len(*p)-1]
	last_worker.index = -1 // for safety
	pool := (*p)[0 : len(*p)-1]
	*p = pool
	return last_worker
}

func (p *Pool) Swap(i, j int) {
	pool := *p
	pool[i], pool[j] = pool[j], pool[i]
	pool[i].index = i
	pool[j].index = j
}

// The worker with the smallest number of pending requests will bubble
// to the top of the heap.
func (p Pool) Less(i, j int) bool {
	return p[i].pending < p[j].pending
}

// Send a Request to a worker
func (b *Balancer) dispatch(req *Request) {
	least_loaded_worker := heap.Pop(&b.pool).(*Worker)
	fmt.Println(fmt.Sprintf("Dispatching request %v to %v", req, least_loaded_worker))
	least_loaded_worker.requests <- *req
	least_loaded_worker.pending++
	heap.Push(&b.pool, least_loaded_worker)
}

// Job is complete, update the pool and the worker to reflect this.
func (b *Balancer) completed(worker *Worker) {
	worker.pending--
	heap.Remove(&b.pool, worker.index)
	heap.Push(&b.pool, worker)
}
