package main

import (
	"sync"
	"time"
)

type Task struct {
	Id             string    `json:"index"`
	N              int       `json:"n"`
	D              float32   `json:"d"`
	N1             float32   `json:"n1"`
	I              float32   `json:"i"`
	TTL            float32   `json:"ttl"`
	Result         float32   `json:"res"`
	Status         string    `json:"status"`
	TimeStart      time.Time `json:"time_start"`
	TimeFinish     time.Time `json:"time_finish"`
	TimeProccesing time.Time `json:"time_proccesing"`
}

type TaskList struct {
	tasks   []*Task
	wgWork  *sync.WaitGroup
	wgCalc  *sync.WaitGroup
	mu      *sync.Mutex
	chtasks chan *Task
}

func NewTaskList() *TaskList {
	var tl *TaskList = &TaskList{}
	tl.chtasks = make(chan *Task, 100)
	tl.mu = &sync.Mutex{}
	tl.wgWork = &sync.WaitGroup{}
	tl.wgCalc = &sync.WaitGroup{}
	tl.tasks = make([]*Task, 0)
	return tl
}
