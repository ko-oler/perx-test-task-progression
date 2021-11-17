package main

import (
	"context"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"time"
)

//Logic
func (tl *TaskList) Work(ctx context.Context) {
	defer tl.wgWork.Done()
	for {
		select {
		case <-ctx.Done():
			tl.wgCalc.Wait()
			log.Println("Done, shutting down the tasks")
			return
		case t := <-tl.chtasks:
			tl.wgCalc.Add(1)
			tl.CalcProgression(t)
			tl.wgCalc.Done()
		}
	}

}

func (tl *TaskList) AddTaskToQueue(t *Task) {
	t.Id = strconv.Itoa(rand.Intn(1000000))
	t.Status = "Pending task"
	t.TimeStart = time.Now()
	tl.mu.Lock()
	tl.tasks = append(tl.tasks, t)
	tl.mu.Unlock()
	tl.chtasks <- t
}

func (tl *TaskList) DeleteTask(t *Task) []*Task {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	for index, task := range tl.tasks {
		if task.Id == t.Id {
			log.Printf("Task deleted by TTL, ID:%s, TTL:%v seconds", t.Id, t.TTL)
			tl.tasks = append(tl.tasks[:index], tl.tasks[index+1:]...)
			break
		}
	}

	return tl.tasks
}

func (tl *TaskList) SortTasks() {
	sort.Slice(tl.tasks, func(p, q int) bool {
		return tl.tasks[p].TimeStart.Before(tl.tasks[q].TimeStart)
	})
}

func (tl *TaskList) CalcProgression(t *Task) {
	// defer tl.wgCalc.Done()
	var res float32 = t.N1
	t.Status = "Processing task"
	t.TimeProccesing = time.Now()
	log.Printf("Started calculation, ID:%s", t.Id)
	for i := 0; i < t.N; i++ {
		res += t.D
		time.Sleep(time.Duration(t.I) * time.Millisecond)
	}
	t.Result = res
	t.Status = "Finished task"
	t.TimeFinish = time.Now()
	log.Printf("Finished task, ID:%s, result:%v ", t.Id, t.Result)
	time.AfterFunc(time.Duration(t.TTL)*time.Second, func() { tl.DeleteTask(t) })
}
