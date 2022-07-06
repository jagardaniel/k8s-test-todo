package main

type InMemoryTodoService struct {
	tasks  map[int]Task
	nextID int
}

func NewInMemoryTodoService() *InMemoryTodoService {
	ts := &InMemoryTodoService{}
	ts.tasks = make(map[int]Task)
	ts.nextID = 1

	return ts
}

// CreateTask creates a new task in the store.
func (i *InMemoryTodoService) CreateTask(title string) int {
	newID := i.nextID
	task := Task{
		ID:        newID,
		Title:     title,
		Completed: false,
	}

	i.tasks[i.nextID] = task
	i.nextID++

	return newID
}

// DeleteTask delets a task with a specific ID. If the id does not exist, return an error
func (i *InMemoryTodoService) DeleteTask(id int) error {
	_, ok := i.tasks[id]
	if ok {
		delete(i.tasks, id)
		return nil
	} else {
		return ErrNotFound
	}
}

func (i *InMemoryTodoService) UpdateTask(id int, title string, completed bool) error {
	_, ok := i.tasks[id]
	if ok {
		updatedTask := Task{
			ID:        id,
			Title:     title,
			Completed: completed,
		}

		i.tasks[id] = updatedTask

		return nil
	} else {
		return ErrNotFound
	}
}

func (i *InMemoryTodoService) GetTasks() []Task {
	allTasks := make([]Task, 0, len(i.tasks))
	for _, task := range i.tasks {
		allTasks = append(allTasks, task)
	}

	return allTasks
}

func (i *InMemoryTodoService) GetTask(id int) (Task, error) {
	task, ok := i.tasks[id]
	if ok {
		return task, nil
	} else {
		return Task{}, ErrNotFound
	}
}
