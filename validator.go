package main

type CreateTaskInput struct {
	Title string `json:"title"`
}

type UpdateTaskInput struct {
	Title string `json:"title,omitempty"`
	Done  bool   `json:"done,omitempty"`
}

func (input *CreateTaskInput) validate() map[string]string {
	if input == nil || input.Title == "" {
		return map[string]string{"title": "Title is required"}
	}
	return nil
}

func (input *UpdateTaskInput) validate() map[string]string {
	if input == nil || (input.Title == "" && !input.Done) {
		return map[string]string{"field": "At least one field (title or done) must be provided"}
	}
	return nil
}
