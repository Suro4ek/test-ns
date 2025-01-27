package entities

import (
	"time"
)

type Task struct {
	id          uint32
	title       string
	description string
	status      string
	createdAt   time.Time
	updatedAt   *time.Time
}

func NewTask(id uint32, title, description, status string, createdAt time.Time, updatedAt *time.Time) Task {
	return Task{
		id:          id,
		title:       title,
		description: description,
		status:      status,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}

func NewTaskCreate(
	title, desscription string, status string, createdAt time.Time,
) Task {
	return Task{
		title:       title,
		description: desscription,
		status:      status,
		createdAt:   createdAt,
	}
}

func (t Task) ID() uint32 {
	return t.id
}

func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}

func (t Task) Status() string {
	return t.status
}

func (t Task) CreatedAt() time.Time {
	return t.createdAt
}

func (t Task) UpdatedAt() *time.Time {
	return t.updatedAt
}

func (t *Task) SetID(id uint32) {
	t.id = id
}

func (t *Task) SetTitle(title string) {
	t.title = title
}

func (t *Task) SetDescription(description string) {
	t.description = description
}

func (t *Task) SetStatus(status string) {
	t.status = status
}

func (t *Task) SetUpdatedAt(updatedAt *time.Time) {
	t.updatedAt = updatedAt
}
