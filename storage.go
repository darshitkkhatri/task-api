package main

type Storage interface {
	Add(title string) (Task, error)
	GetAll() ([]Task, error)
	GetByID(id int) (Task, bool, error)
	Update(id int, title string, done bool) (Task, bool, error)
	Delete(id int) (bool, error)
}
