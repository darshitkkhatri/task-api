package main

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
	"time"
)

type Store struct {
	mu       sync.Mutex
	filePath string
}

func NewStore(filePath string) *Store {
	return &Store{filePath: filePath}
}

var _ Storage = (*Store)(nil)

func (s *Store) load() (map[int]Task, int, error) {
	tasks := make(map[int]Task)

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return tasks, 1, nil
		}
		return nil, 0, err
	}

	if len(data) == 0 {
		return tasks, 1, nil
	}

	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, 0, err
	}

	nextID := 1
	for id := range tasks {
		if id >= nextID {
			nextID = id + 1
		}
	}
	return tasks, nextID, nil
}

func (s *Store) save(tasks map[int]Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.filePath, data, 0644)
}

func (s *Store) Add(title string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, nextID, err := s.load()
	if err != nil {
		return Task{}, err
	}

	task := Task{
		ID:        nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	tasks[nextID] = task
	return task, s.save(tasks)
}

func (s *Store) GetAll() ([]Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, _, err := s.load()
	if err != nil {
		return nil, err
	}

	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, task)
	}
	return result, nil
}

func (s *Store) GetByID(id int) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, _, err := s.load()
	if err != nil {
		return Task{}, false, err
	}

	task, ok := tasks[id]
	return task, ok, nil
}

func (s *Store) Update(id int, title string, done bool) (Task, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, _, err := s.load()
	if err != nil {
		return Task{}, false, err
	}

	task, ok := tasks[id]
	if !ok {
		return Task{}, false, nil
	}

	task.Title = title
	task.Done = done
	tasks[id] = task
	return task, true, s.save(tasks)
}

func (s *Store) Delete(id int) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tasks, _, err := s.load()
	if err != nil {
		return false, err
	}

	_, ok := tasks[id]
	if !ok {
		return false, nil
	}

	delete(tasks, id)
	return true, s.save(tasks)
}

func parseID(s string) (int, bool) {
	id, err := strconv.Atoi(s)
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}
