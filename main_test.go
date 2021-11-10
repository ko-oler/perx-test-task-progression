package main

import (
	"log"
	"sync"
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
	var mu = &sync.Mutex{}
	var wg = &sync.WaitGroup{}
	var pending, finished = "Pending task", "Finished task"

	for w := 0; w < totalWorker; w++ {
		go Work(mu, wg)
	}

	for _, pair := range test {
		AddTaskToQueue(pair.task)

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
				"expected", pair.res,
				"got", pair.task.Result,
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

	for _, v := range tasks {
		if v.Status != finished {
			t.Error(
				"For task ID", v.Id,
				"expected", finished,
				"got", v.Status,
			)
		}
	}
}

func TestCalcProgression(t *testing.T) {
	var mu = &sync.Mutex{}
	var wg = &sync.WaitGroup{}
	for _, pair := range test {
		wg.Add(1)
		CalcProgression(pair.task, mu, wg)
		if pair.task.Result != pair.res {
			t.Error(
				"For task ID", pair.task.Id,
				"expected", pair.res,
				"got", pair.task.Result,
			)

		}
	}
}
