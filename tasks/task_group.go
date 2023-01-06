package tasks

import "sync"

type TaskGroup struct {
	Name  string
	Tasks []*Task
	Mutex *sync.Mutex
}

// NewTaskGroup Create and returns a new task group
func NewTaskGroup(name string) TaskGroup {
	g := TaskGroup{
		Name:  name,
		Tasks: []*Task{},
		Mutex: &sync.Mutex{},
	}

	return g
}

// AddTask Adds a task to the group
func (g *TaskGroup) AddTask(t *Task) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	t.Group = g
	g.Tasks = append(g.Tasks, t)
}

// RemoveTask Removes a task from the group
func (g *TaskGroup) RemoveTask(t *Task) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	t.Group = nil

	for i, task := range g.Tasks {
		if task == t {
			g.Tasks = append(g.Tasks[:i], g.Tasks[i+1:]...)
		}
	}
}
