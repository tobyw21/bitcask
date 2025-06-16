package fs

import "container/heap"

/*
	pq.go contains a priority queue used for vfd lru

*/

type OpenedVfdItem struct {
	vfd_id uint8
	// use uint8 here because we keep the queue relatively small
	priority uint8
	index    int8
}

type VfdPriorityQueue []*OpenedVfdItem

func (pq VfdPriorityQueue) Len() int { return len(pq) }

func (pq VfdPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = int8(i)
	pq[j].index = int8(j)
}

func (pq VfdPriorityQueue) Less(i, j int) bool {
	// pop will give us the lowest priority aka, the usage count of fd
	return pq[i].priority < pq[j].priority
}

func (pq *VfdPriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*OpenedVfdItem)
	item.index = int8(n)
	*pq = append(*pq, item)
}

func (pq *VfdPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*pq = old[0 : n-1]
	return item.vfd_id
}




func (pq *VfdPriorityQueue) Update(fd uint8, priority uint8, index int8) {
	item := &OpenedVfdItem{
		vfd_id:   fd,
		priority: priority,
		index:    int8(index),
	}

	heap.Fix(pq, int(item.index))
}
