package gox

//export
func RunFunc(processCount int, funcs ...func()) {
	tasks := make([]Task, len(funcs))
	for i, funcItem := range funcs {
		tasks[i] = SimpleTask(funcItem)
	}
	RunTask(processCount, tasks...)
}

//export
func RunFuncWithCallback(processCount int, callback TaskCallback, funcs ...func() (interface{}, error)) {
	tasks := make([]Task, len(funcs))
	for i, funcItem := range funcs {
		tasks[i] = TaskFunc(funcItem)
	}
	RunTaskWithCallback(processCount, callback, tasks...)
}
