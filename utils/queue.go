package utils

import (
	"sync"
)

type Queue struct {
	items []interface{} // 队列的数据存储
	lock  sync.Mutex    // 用于保护共享资源的互斥锁
	cond  *sync.Cond    // 条件变量，处理阻塞和唤醒
}

func NewQueue() *Queue {
	q := &Queue{}
	q.cond = sync.NewCond(&q.lock) // 将条件变量与互斥锁关联
	return q
}

func (q *Queue) Append(value interface{}) {
	q.lock.Lock()         // 加锁以保护共享资源
	defer q.lock.Unlock() // 确保函数退出时解锁

	q.items = append(q.items, value) // 将数据添加到队列
	q.cond.Signal()                  // 唤醒一个等待的 goroutine
}

func (q *Queue) Pop() interface{} {
	q.lock.Lock()         // 加锁以保护共享资源
	defer q.lock.Unlock() // 确保函数退出时解锁

	// 当队列为空时，进入等待状态，直到有新数据入队
	for len(q.items) == 0 {
		q.cond.Wait() // 阻塞当前 goroutine，等待数据
	}

	value := q.items[0]   // 读取队列中的第一个元素
	q.items = q.items[1:] // 从队列中移除第一个元素
	return value          // 返回出队的数据
}
