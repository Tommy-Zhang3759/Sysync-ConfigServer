package clientManage

import (
	"container/list"
	"time"
)

type Task interface {
	Run() error
	Stop() error
}

type Schedule struct {
	ExecTime time.Time
	Do       func() error
}

func (s *Schedule) Run() error {
	return s.Do()
}
func (s *Schedule) Stop() error {
	return nil
}

type TimeWheel struct {
	slots     []list.List
	slotCount int
	tickRate  time.Duration
	current   int
}

func NewTimeWheel(slotCount int, tickRate time.Duration) *TimeWheel {
	slots := make([]list.List, slotCount)
	return &TimeWheel{
		slots:     slots,
		slotCount: slotCount,
		tickRate:  tickRate,
	}
}

func (tw *TimeWheel) AddTask(task *Schedule) {
	delay := task.ExecTime.Sub(time.Now())
	if delay < 0 {
		delay = 0
	}
	slotIndex := int(delay/tw.tickRate) % tw.slotCount
	tw.slots[slotIndex].PushBack(task)
}

func (tw *TimeWheel) Start() {
	for {
		time.Sleep(tw.tickRate)
		tw.current = (tw.current + 1) % tw.slotCount
		slot := tw.slots[tw.current]
		for e := slot.Front(); e != nil; e = e.Next() {
			task := e.Value.(*Schedule)
			if time.Now().After(task.ExecTime) {
				task.Run()
				slot.Remove(e)
			}
		}
	}
}

/*
func main() {
	tw := NewTimeWheel(10, 1*time.Second)

	tw.AddTask(&Schedule{
		ExecTime: time.Now().Add(3 * time.Second),
		action: func() {
			fmt.Println("Schedule executed after 3 seconds")
		},
	})

	tw.AddTask(&Schedule{
		ExecTime: time.Now().Add(5 * time.Second),
		action: func() {
			fmt.Println("Schedule executed after 5 seconds")
		},
	})

	go tw.Start()

	// Keep the main function alive to observe the time wheel in action
	time.Sleep(10 * time.Second)
}
*/
