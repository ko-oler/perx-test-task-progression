package main

import (
	"context"
	"log"
	"testing"
	"time"
)

type testpair struct {
	task *Task
	res  float32
}

var test = []testpair{
	testpair{
		task: &Task{
			Id:             "1",
			N:              10,
			D:              1,
			N1:             0,
			I:              10,
			TTL:            2,
			Result:         0,
			Status:         "",
			TimeStart:      time.Time{},
			TimeFinish:     time.Time{},
			TimeProccesing: time.Time{},
		},
		res: 10,
	},
	testpair{
		task: &Task{
			Id:             "2",
			N:              100,
			D:              1,
			N1:             0,
			I:              10,
			TTL:            2,
			Result:         0,
			Status:         "",
			TimeStart:      time.Time{},
			TimeFinish:     time.Time{},
			TimeProccesing: time.Time{},
		},
		res: 100,
	},
	testpair{
		task: &Task{
			Id:             "3",
			N:              20,
			D:              2,
			N1:             0,
			I:              10,
			TTL:            2,
			Result:         0,
			Status:         "",
			TimeStart:      time.Time{},
			TimeFinish:     time.Time{},
			TimeProccesing: time.Time{},
		},
		res: 40,
	},
}

func TestWork(t *testing.T) {

	var totalWorker = 2
	tl := NewTaskList()
	ctx, cancel := context.WithCancel(context.Background())
	var pending, finished = "Pending task", "Finished task"

	for w := 0; w < totalWorker; w++ {
		tl.wgWork.Add(1)
		go tl.Work(ctx)
	}

	for _, pair := range test {
		tl.AddTaskToQueue(pair.task)

		if pair.task.Status != pending {
			t.Error(
				"For task ID", pair.task.Id,
				"expected", pending,
				"got", pair.task.Status,
			)
		}
		if pair.task.Id == "1" || pair.task.Id == "2" || pair.task.Id == "3" {
			t.Error(
				"For task ID", pair.task.Id,
				"expected ID to be generated",
				"got", pair.task.Id,
			)
		}

		if pair.task.TimeStart.IsZero() {
			t.Error(
				"For task ID", pair.task.Id,
				"expected smth about this", time.Now(),
				"got", pair.task.TimeStart,
			)
		}
		log.Println(pair.task.Id)

	}

	time.Sleep(2 * time.Second)

	for _, v := range tl.tasks {
		if v.Status != finished {
			t.Error(
				"For task ID", v.Id,
				"expected", finished,
				"got", v.Status,
			)
		}
	}

	for _, v := range test {
		if v.task.Result != v.res {
			t.Error(
				"For task ID", v.task.Id,
				"expected", v.res,
				"got", v.task.Result,
			)
		}
	}

	cancel()
	tl.wgWork.Wait()

}

func TestCalcProgression(t *testing.T) {
	tl := NewTaskList()
	ctx, cancel := context.WithCancel(context.Background())
	for _, pair := range test {
		tl.wgCalc.Add(1)
		tl.CalcProgression(pair.task, ctx)
		if pair.task.Result != pair.res {
			t.Error(
				"For task ID", pair.task.Id,
				"expected", pair.res,
				"got", pair.task.Result,
			)

		}
	}
	cancel()
}
