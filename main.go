// Сервис, вычисляющий арифметическую прогрессию в очереди.
// Задачи поступают в очередь, из очереди поступают на выполнение, выполняются до получения результата, после чего из очереди выбирается следующая задача.
// Параллельно может выполняться N задач. Количество N передается через параметры командной строки сервиса.

// Параметры задачи:
// n - количество элементов целочисленное (целочисленное)
// d - дельта между элементами последовательности (вещественное)
// n1 - Стартовое значение (вещественное)
// I - интервал в секундах между итерациями (вещественное)
// TTL - время хранения результата в секундах (вещественное)

// При запуске стартует HTTP-сервер, у сервера есть два endpointa:
// Постановка задачи в очередь.
// Получение отсортированного списка задач и статусы выполнения этих задач.
// Отработанные задачи стираются после завершения TTL.

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const totalWorker int = 2

type Task struct {
	Id             string  `json:"index"`
	N              int     `json:"n"`
	D              float32 `json:"d"`
	N1             float32 `json:"n1"`
	I              float32 `json:"i"`
	TTL            float32 `json:"ttl"`
	Result         float32 `json:"res"`
	Status         string  `json:"status"`
	TimeStart      string  `json:"ts"`
	TimeFinish     string  `json:"tf"`
	TimeProccesing string  `json:"tp"`
}

var chtasks = make(chan *Task, 100)
var tasks []*Task

func getlistoftasks(w http.ResponseWriter, r *http.Request) {

	sort.Slice(tasks, func(p, q int) bool {
		return tasks[p].TimeStart < tasks[q].TimeStart
	})
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func Work(mu *sync.Mutex) {
	for t := range chtasks {
		//wg.Add(1)
		calcprogression(t, mu)
		//wg.Done()
	}
}

func addtask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var task *Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	task.Id = strconv.Itoa(rand.Intn(1000000))
	task.Status = "Pending task"
	task.TimeStart = time.Now().Format(time.RFC850)
	tasks = append(tasks, task)
	json.NewEncoder(w).Encode(task)
	chtasks <- task
}

func calcprogression(t *Task, mu *sync.Mutex) {
	var res float32 = t.N1
	fmt.Println("Started calculation, ID:", t.Id)
	t.Status = "Processing task"
	t.TimeProccesing = time.Now().Format(time.RFC850)
	for i := 0; i < t.N; i++ {
		res += t.D
		time.Sleep(time.Duration(t.I) * time.Millisecond)
	}
	t.Result = res
	fmt.Println("Finished task, ID:", t.Id)
	t.Status = "Finished task"
	t.TimeFinish = time.Now().Format(time.RFC850)
	time.AfterFunc(time.Duration(t.TTL)*time.Second, func() { deletetask(t, mu) })
}

func deletetask(t *Task, mu *sync.Mutex) []*Task {
	mu.Lock()
	for index, task := range tasks {
		if task.Id == t.Id {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}
	mu.Unlock()
	return tasks
}

func main() {

	//wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	r := mux.NewRouter()
	r.HandleFunc("/List", getlistoftasks).Methods("GET")
	r.HandleFunc("/Add", addtask).Methods("POST")
	fmt.Println("Server started")
	for t := 0; t < totalWorker; t++ {
		// wg.Add(1)
		go Work(mu)
	}
	log.Fatal(http.ListenAndServe(":80", r))
	//wg.Wait()
}
