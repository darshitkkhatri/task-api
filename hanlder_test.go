package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupHandler(t *testing.T) *Handler {
	return NewHandler(newMockStorage())
}

func makeRequest(t *testing.T, handler http.HandlerFunc, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var req *http.Request

	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatal("failed to marshal body:", err)
		}
		t.Logf("request body: %s", string(data))
		req = httptest.NewRequest(method, path, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	rr := httptest.NewRecorder()
	handler(rr, req)

	t.Logf("response status: %d", rr.Code)
	t.Logf("response body: %s", rr.Body.String())

	return rr
}

func TestHandler_ListTasks_Empty(t *testing.T) {
	h := setupHandler(t)

	rr := makeRequest(t, h.listTasks, "GET", "/tasks", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var tasks []Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 0 {
		t.Errorf("expected empty list, got %d tasks", len(tasks))
	}
}

func TestHandler_CreateTask(t *testing.T) {
	h := setupHandler(t)

	rr := makeRequest(t, h.createTask, "POST", "/tasks", map[string]string{
		"title": "Learn Go",
	})

	if rr.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", rr.Code)
	}

	var task Task
	json.NewDecoder(rr.Body).Decode(&task)
	if task.Title != "Learn Go" {
		t.Errorf("expected title 'Learn Go', got '%s'", task.Title)
	}
	if task.ID == 0 {
		t.Error("expected task to have an ID")
	}
}

func TestHandler_CreateTask_InvalidBody(t *testing.T) {
	h := setupHandler(t)

	rr := makeRequest(t, h.createTask, "POST", "/tasks", map[string]string{
		"title": "",
	})

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rr.Code)
	}

	var apiErr APIError
	json.NewDecoder(rr.Body).Decode(&apiErr)
	if apiErr.Fields["title"] == "" {
		t.Error("expected title field error")
	}
}

func TestHandler_CreateTask_TitleTooLong(t *testing.T) {
	h := setupHandler(t)

	// use real characters instead of null bytes
	longTitle := strings.Repeat("a", 200)
	rr := makeRequest(t, h.createTask, "POST", "/tasks", map[string]string{
		"title": longTitle,
	})

	if rr.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", rr.Code)
	}
}

func TestHandler_ListTasks_AfterCreate(t *testing.T) {
	h := setupHandler(t)

	makeRequest(t, h.createTask, "POST", "/tasks", map[string]string{"title": "Task 1"})
	makeRequest(t, h.createTask, "POST", "/tasks", map[string]string{"title": "Task 2"})

	rr := makeRequest(t, h.listTasks, "GET", "/tasks", nil)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var tasks []Task
	json.NewDecoder(rr.Body).Decode(&tasks)
	if len(tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasks))
	}
}
