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
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
)

func main() {
	tl := NewTaskList()
	ctx, cancel := context.WithCancel(context.Background())
	var totalWorker = 2

	flag.IntVar(&totalWorker, "tw", 2, "number of proccesing task")
	flag.Parse()
	log.Println("Number of proccesing task: ", totalWorker)

	if totalWorker == 0 {
		log.Fatalf("Zero workers")
	}

	r := mux.NewRouter()
	httpServer := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	r.HandleFunc("/Add", AddHandler(tl)).Methods("POST")
	r.HandleFunc("/List", GetListHandler(tl)).Methods("GET")
	log.Println("Server started on port :80")

	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatalf("HTTP server closed with: %v\n", err)
			}
			log.Printf("HTTP server shut down")
		}
	}()

	// Start workerpool
	for t := 0; t < totalWorker; t++ {
		tl.wgWork.Add(1)
		go tl.Work(ctx)
	}
	<-termChan
	log.Println("SIGTERM received. Shutdown process initiated")
	close(tl.chtasks)
	// Shutdown the HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	// Cancel the context
	cancel()
	// Wait to finish
	log.Println("waiting current tasks to finish...")
	tl.wgWork.Wait()
	log.Println("Done. returning.")
}
