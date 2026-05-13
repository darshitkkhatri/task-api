package main

import (
	"os"
	"testing"
)

func setupStore(t *testing.T) *Store {
	f, err := os.CreateTemp("", "tasks_test_*.json")
	if err != nil {
		t.Fatal("failed to create temp file:", err)
	}
	f.Close()

	t.Cleanup(func() {
		os.Remove(f.Name())
	})

	return NewStore(f.Name())
}

func TestStore_Add(t *testing.T) {
	store := setupStore(t)

	task, err := store.Add("Learn Go")
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}
	if task.Title != "Learn Go" {
		t.Errorf("expected title 'Learn Go', got '%s'", task.Title)
	}
	if task.Done != false {
		t.Error("expected Done to be false")
	}
}

func TestStore_Add_MultipleItems(t *testing.T) {
	store := setupStore(t)

	first, _ := store.Add("Task 1")
	second, _ := store.Add("Task 2")

	if first.ID == second.ID {
		t.Error("expected different IDs for different tasks")
	}
	if second.ID != first.ID+1 {
		t.Errorf("expected sequential IDs, got %d and %d", first.ID, second.ID)
	}
}

func TestStore_GetAll(t *testing.T) {
	store := setupStore(t)

	tasks, err := store.GetAll()
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks, got %d", len(tasks))
	}

	store.Add("Task 1")
	store.Add("Task 2")

	tasks, err = store.GetAll()
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}

func TestStore_GetByID(t *testing.T) {
	store := setupStore(t)

	added, _ := store.Add("Learn Go")

	task, ok, err := store.GetByID(added.ID)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if !ok {
		t.Fatal("expected task to exist")
	}
	if task.Title != "Learn Go" {
		t.Errorf("expected title 'Learn Go', got '%s'", task.Title)
	}

	_, ok, err = store.GetByID(999)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if ok {
		t.Error("expected task to not exist")
	}
}

func TestStore_Update(t *testing.T) {
	store := setupStore(t)

	added, _ := store.Add("Learn Go")

	updated, ok, err := store.Update(added.ID, "Learn Go well", true)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if !ok {
		t.Fatal("expected task to exist")
	}
	if updated.Title != "Learn Go well" {
		t.Errorf("expected title 'Learn Go well', got '%s'", updated.Title)
	}
	if !updated.Done {
		t.Error("expected Done to be true")
	}

	_, ok, err = store.Update(999, "nothing", false)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if ok {
		t.Error("expected task to not exist")
	}
}

func TestStore_Delete(t *testing.T) {
	store := setupStore(t)

	added, _ := store.Add("Learn Go")

	ok, err := store.Delete(added.ID)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if !ok {
		t.Fatal("expected task to be deleted")
	}

	_, ok, _ = store.GetByID(added.ID)
	if ok {
		t.Error("expected task to be gone after delete")
	}

	ok, err = store.Delete(999)
	if err != nil {
		t.Fatal("expected no error, got:", err)
	}
	if ok {
		t.Error("expected false for non existing task")
	}
}
