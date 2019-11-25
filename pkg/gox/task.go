package gox

import (
	"fmt"
	"sync"
)

type Task interface {
	Run() (interface{}, error)
}

type SimpleTask func()

func (st SimpleTask) Run() (interface{}, error) {
	st()
	return nil, nil
}

type TaskFunc func() (interface{}, error)

func (tf TaskFunc) Run() (interface{}, error) {
	return tf()
}

type TaskCallback func(interface{}, error)

func startTaskProcessor(processCount int, wg *sync.WaitGroup, callback TaskCallback) chan Task {
	taskQueue := make(chan Task)
	for i := 0; i < processCount; i++ {
		go goTask(taskQueue, wg, callback)
	}

	return taskQueue
}

func goTask(taskQueue chan Task, wg *sync.WaitGroup, callback TaskCallback) {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	wg.Add(1)
	for {
		task, ok := <-taskQueue
		if !ok {
			break
		}
		//signal := NewStateSignal()
		runTask(task, wg, callback)
		//signal.Wait(1)
	}
}

func runTask(task Task, wg *sync.WaitGroup, callback TaskCallback) {
	defer func() {
		wg.Done()
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	wg.Add(1)
	result, err := task.Run()
	if callback != nil {
		callback(result, err)
	}
}

func RunTaskWithCallback(processCount int, callback TaskCallback, tasks ...Task) {
	var wg sync.WaitGroup
	taskQueue := startTaskProcessor(processCount, &wg, callback)
	for _, t := range tasks {
		taskQueue <- t
	}
	close(taskQueue)

	wg.Wait()
}

func RunTask(processCount int, tasks ...Task) {
	RunTaskWithCallback(processCount, nil, tasks...)
}
