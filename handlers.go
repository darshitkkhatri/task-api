package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	storage Storage
}

func NewHandler(storage Storage) *Handler {
	return &Handler{storage: storage}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.storage.GetAll()
	if err != nil {
		internalError("failed to fetch tasks").respond(w)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *Handler) createTask(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil || r.ContentLength == 0 {
		badRequest("request body is required").respond(w)
		return
	}

	var input CreateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest("invalid JSON body: " + err.Error()).respond(w)
		return
	}

	if errs := input.validate(); len(errs) > 0 {
		validationError(errs).respond(w)
		return
	}

	task, err := h.storage.Add(input.Title)
	if err != nil {
		internalError("failed to create task").respond(w)
		return
	}
	writeJSON(w, http.StatusCreated, task)
}

func (h *Handler) getTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		badRequest("invalid task ID").respond(w)
		return
	}

	task, ok, err := h.storage.GetByID(id)
	if err != nil {
		internalError("failed to fetch task").respond(w)
		return
	}
	if !ok {
		notFound("task not found").respond(w)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) updateTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		badRequest("invalid task ID").respond(w)
		return
	}

	var input UpdateTaskInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		badRequest("invalid JSON body: " + err.Error()).respond(w)
		return
	}

	if errs := input.validate(); len(errs) > 0 {
		validationError(errs).respond(w)
		return
	}

	task, ok, err := h.storage.Update(id, input.Title, input.Done)
	if err != nil {
		internalError("failed to update task").respond(w)
		return
	}
	if !ok {
		notFound("task not found").respond(w)
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *Handler) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r.PathValue("id"))
	if !ok {
		badRequest("invalid task ID").respond(w)
		return
	}

	ok, err := h.storage.Delete(id)
	if err != nil {
		internalError("failed to delete task").respond(w)
		return
	}
	if !ok {
		notFound("task not found").respond(w)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "task deleted"})
}

func (h *Handler) healthCheck(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": "1.0",
	})
}
