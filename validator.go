package main

import "strings"

type CreateTaskInput struct {
	Title string `json:"title"`
}

type UpdateTaskInput struct {
	Title string `json:"title,omitempty"`
	Done  bool   `json:"done,omitempty"`
}

func (input *CreateTaskInput) validate() map[string]string {
	errs := make(map[string]string)

	title := strings.TrimSpace(input.Title)

	if title == "" {
		errs["title"] = "Title is required"
	} else if len([]rune(title)) > 100 { // ← count characters not bytes
		errs["title"] = "title must be under 100 characters"
	}

	return errs

}

func (input *UpdateTaskInput) validate() map[string]string {
	errs := make(map[string]string)

	title := strings.TrimSpace(input.Title)

	if title == "" {
		errs["title"] = "Title is required"
	} else if len([]rune(title)) > 100 { // ← count characters not bytes
		errs["title"] = "title must be under 100 characters"
	}

	return errs
}
