package main

type mockStorage struct {
	tasks  map[int]Task
	nextID int
	err    error
}

func newMockStorage() *mockStorage {
	return &mockStorage{
		tasks:  make(map[int]Task),
		nextID: 1,
	}
}

func (m *mockStorage) Add(title string) (Task, error) {
	if m.err != nil {
		return Task{}, m.err
	}
	task := Task{ID: m.nextID, Title: title, Done: false}
	m.tasks[m.nextID] = task
	m.nextID++
	return task, nil
}

func (m *mockStorage) GetAll() ([]Task, error) {
	if m.err != nil {
		return nil, m.err
	}
	result := make([]Task, 0, len(m.tasks))
	for _, t := range m.tasks {
		result = append(result, t)
	}
	return result, nil
}

func (m *mockStorage) GetByID(id int) (Task, bool, error) {
	if m.err != nil {
		return Task{}, false, m.err
	}
	task, ok := m.tasks[id]
	return task, ok, nil
}

func (m *mockStorage) Update(id int, title string, done bool) (Task, bool, error) {
	if m.err != nil {
		return Task{}, false, m.err
	}
	task, ok := m.tasks[id]
	if !ok {
		return Task{}, false, nil
	}
	task.Title = title
	task.Done = done
	m.tasks[id] = task
	return task, true, nil
}

func (m *mockStorage) Delete(id int) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	_, ok := m.tasks[id]
	if !ok {
		return false, nil
	}
	delete(m.tasks, id)
	return true, nil
}
