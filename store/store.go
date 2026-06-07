package store

import (
	"fmt"
	"sync"
	"task-api-comparison/model"
)

var (
	mu     sync.RWMutex
	tasks  = map[string]*model.Task{
		"t_001": {ID: "t_001", Title: "Первая задача", Description: "Учебный пример", Done: false},
		"t_002": {ID: "t_002", Title: "Вторая задача", Description: "Проверка API", Done: true},
	}
	nextID = 3
)

func All() []*model.Task {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]*model.Task, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, t)
	}
	return result
}

func ByID(id string) (*model.Task, bool) {
	mu.RLock()
	defer mu.RUnlock()
	t, ok := tasks[id]
	return t, ok
}

func Create(title, desc string) *model.Task {
	mu.Lock()
	defer mu.Unlock()
	id := fmt.Sprintf("t_%03d", nextID)
	nextID++
	t := &model.Task{ID: id, Title: title, Description: desc, Done: false}
	tasks[id] = t
	return t
}

func UpdateDone(id string, done bool) (*model.Task, bool) {
	mu.Lock()
	defer mu.Unlock()
	t, ok := tasks[id]
	if !ok {
		return nil, false
	}
	t.Done = done
	return t, true
}