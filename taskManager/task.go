package TaskManager

import (
	"container/list"
	"fmt"
	"time"
)

type Task interface {
	Run() error
	Stop() error
}

type TimeWhellMember struct {
	ExecTime time.Time
	task     Task
}

func (s *TimeWhellMember) Run() error {
	return s.task.Run()
}

func (s *TimeWhellMember) Stop() error {
	return s.task.Stop()
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

func (tw *TimeWheel) AddTask(task *TimeWhellMember) error {
	delay := time.Until(task.ExecTime)
	if delay < 0 {
		return fmt.Errorf("error: trying to add a expired task into time whell")
	}
	slotIndex := int(delay/tw.tickRate) % tw.slotCount
	tw.slots[slotIndex].PushBack(task)
	return nil
}

func NewTimeWheelMember(exeTime time.Time, exeFunc Task) *TimeWhellMember {
	return &TimeWhellMember{
		ExecTime: exeTime,
		task:     exeFunc,
	}
}

func (tw *TimeWheel) Start() {
	for {
		time.Sleep(tw.tickRate)
		tw.current = (tw.current + 1) % tw.slotCount
		slot := tw.slots[tw.current]
		for e := slot.Front(); e != nil; e = e.Next() {
			go func(task *TimeWhellMember) {
				if time.Now().After(task.ExecTime) {
					task.Run()
					slot.Remove(e)
				}
			}(e.Value.(*TimeWhellMember))
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
