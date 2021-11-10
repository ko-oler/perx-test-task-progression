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
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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

var totalWorker int = 2
var chtasks = make(chan *Task, 100)
var tasks []*Task
var mu = &sync.Mutex{}
var wg = &sync.WaitGroup{}

//Handlers
func GetListHandler(w http.ResponseWriter, r *http.Request) {
	SortTasks()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	w.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	AddTaskToQueue(&task)
	json.NewEncoder(w).Encode(task)

}

//Logic
func Work(mu *sync.Mutex, wg *sync.WaitGroup) {
	for t := range chtasks {
		wg.Add(1)
		CalcProgression(t, mu, wg)
	}
}

func AddTaskToQueue(t *Task) {
	t.Id = strconv.Itoa(rand.Intn(1000000))
	t.Status = "Pending task"
	t.TimeStart = time.Now()
	mu.Lock()
	tasks = append(tasks, t)
	mu.Unlock()
	chtasks <- t
}

func DeleteTask(t *Task, mu *sync.Mutex) []*Task {
	mu.Lock()
	defer mu.Unlock()
	for index, task := range tasks {
		if task.Id == t.Id {
			log.Printf("Task deleted by TTL, ID:%s, TTL:%v seconds", t.Id, t.TTL)
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}
	}

	return tasks
}

func SortTasks() {
	sort.Slice(tasks, func(p, q int) bool {
		return tasks[p].TimeStart.Before(tasks[q].TimeStart)
	})
}

func CalcProgression(t *Task, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
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
	time.AfterFunc(time.Duration(t.TTL)*time.Second, func() { DeleteTask(t, mu) })
}

func main() {

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)

	flag.IntVar(&totalWorker, "tw", 2, "количество одновременно выполняемых задач")
	flag.Parse()
	log.Println("Number of proccesing task: ", totalWorker)

	r := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	r.HandleFunc("/Add", AddHandler).Methods("POST")
	r.HandleFunc("/List", GetListHandler).Methods("GET")

	// Запуск воркеров
	for t := 0; t < totalWorker; t++ {
		go Work(mu, wg)
	}

	log.Println("Server started on port :80")

	go func() {
		<-termChan // Blocks here until interrupted
		log.Print("SIGTERM received. Shutdown process initiated\n")
		httpServer.Shutdown(context.Background())
	}()

	// Blocking
	if err := httpServer.ListenAndServe(); err != nil {
		if err.Error() != "http: Server closed" {
			log.Printf("HTTP server closed with: %v\n", err)
		}
		log.Printf("HTTP server shut down")
	}
	log.Println("waiting for running jobs to finish")
	wg.Wait()
	log.Println("jobs finished. exiting")

}
