package tasks

import (
	"Mystery/automation"
	"log"
	"sync"
)

type Manager struct {
	TaskMutex      *sync.Mutex
	TaskGroupMutex *sync.Mutex
	TaskWaitGroup  *sync.WaitGroup
	Tasks          map[*Task]struct{}
	TaskGroups     map[*TaskGroup]struct{}
	Automations    []*automation.Automation
}

// NewManager Creates a new Task manager object
func NewManager() Manager {
	return Manager{
		TaskMutex:      &sync.Mutex{},
		TaskGroupMutex: &sync.Mutex{},
		TaskWaitGroup:  &sync.WaitGroup{},
		Tasks:          map[*Task]struct{}{},
		TaskGroups:     map[*TaskGroup]struct{}{},
		Automations:    []*automation.Automation{},
	}
}

// AddTask Adds a task to the manager
func (m *Manager) AddTask(t *Task) {
	m.TaskMutex.Lock()
	defer m.TaskMutex.Unlock()

	// Check if the task is already in the manager
	if _, ok := m.Tasks[t]; ok {
		return
	}

	t.TaskManager = m
	m.Tasks[t] = struct{}{}
}

// RemoveTask Removes a task from the manager
func (m *Manager) RemoveTask(t *Task) {
	m.TaskMutex.Lock()
	defer m.TaskMutex.Unlock()

	t.TaskManager = nil
	delete(m.Tasks, t)
}

// AddTaskGroup Adds a task group to the manager
func (m *Manager) AddTaskGroup(g *TaskGroup) {
	m.TaskGroupMutex.Lock()
	defer m.TaskGroupMutex.Unlock()

	// Don't allow the same group to be added twice
	if _, ok := m.TaskGroups[g]; ok {
		return
	}

	m.TaskGroups[g] = struct{}{}

	for _, task := range g.Tasks {
		m.AddTask(task)
	}
}

// RemoveTaskGroup Removes a task group from the manager. This also stops and gets rid of all tasks.
func (m *Manager) RemoveTaskGroup(g *TaskGroup) {
	m.TaskGroupMutex.Lock()
	defer m.TaskGroupMutex.Unlock()

	// Stop and remove all tasks that are in the group
	for i := len(g.Tasks) - 1; i >= 0; i-- {
		task := g.Tasks[i]

		switch t := interface{}(task).(type) {
		case TaskRunner:
			m.StopTask(t)
		}

		g.RemoveTask(task)
		m.RemoveTask(task)
	}

	delete(m.TaskGroups, g)
}

// StartTask Starts a specific task
func (m *Manager) StartTask(t TaskRunner) {
	t.Start()
}

// StopTask Stops a specific task
func (m *Manager) StopTask(t TaskRunner) {
	t.Stop()
}

// WaitForAllTasks Waits for every task to call wg.Done()
func (m *Manager) WaitForAllTasks() {
	m.TaskWaitGroup.Wait()
}

// HandleNewAutomationProducts Handles incoming products from automation
func (m *Manager) HandleNewAutomationProducts() {
	log.Println("Starting automation product handler.")

	go func() {
		for {
			select {
			case product := <-automation.ZephyrMonitorChannel:
				for _, auto := range m.Automations {
					if !auto.IsPriceMatch(product) || !auto.IsWebsiteMatch(product) || !auto.IsProductMatch(product) {
						continue
					}

					variants := auto.GetMatchingSizeVariants(product)

					if len(variants) == 0 {
						continue
					}

					// TODO: Check for existing automation group with the exact same variants.
					auto.SendWebhook(true, product)
				}
			}
		}
	}()
}
