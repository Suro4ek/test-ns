package repo

import (
	"context"
	"test-ns/internal/entities"
	"time"
)

type taskGORM struct {
	Id          uint32 `gorm:"primaryKey"`
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time `gorm:"autoUpdateTime:false"`
}

type TaskRepo struct {
	db GSQL
}

func NewTaskRepo(db GSQL) TaskRepo {
	db.AutoMigrate(&taskGORM{})
	return TaskRepo{db: db}
}

func (r TaskRepo) Create(ctx context.Context, task entities.Task) (entities.Task, error) {
	t := taskGORM{
		Id:          task.ID(),
		Title:       task.Title(),
		Description: task.Description(),
		Status:      task.Status(),
		CreatedAt:   task.CreatedAt(),
		UpdatedAt:   task.UpdatedAt(),
	}
	if err := r.db.Create(ctx, &t); err != nil {
		return entities.Task{}, entities.ErrorTaskCreate
	}
	return entities.NewTask(
		t.Id,
		t.Title,
		t.Description,
		t.Status,
		t.CreatedAt,
		t.UpdatedAt,
	), nil
}

func (r TaskRepo) Update(ctx context.Context, task entities.Task) (entities.Task, error) {
	updates := map[string]interface{}{
		"title":       task.Title(),
		"description": task.Description(),
		"status":      task.Status(),
		"created_at":  task.CreatedAt(),
		"updated_at":  task.UpdatedAt(),
	}
	if err := r.db.Update(ctx, &updates, &taskGORM{Id: task.ID()}); err != nil {
		return entities.Task{}, entities.ErrorTaskCreate
	}
	return entities.NewTask(
		task.ID(),
		updates["title"].(string),
		updates["description"].(string),
		updates["status"].(string),
		updates["created_at"].(time.Time),
		updates["updated_at"].(*time.Time),
	), nil
}

func (r TaskRepo) GetTasks(ctx context.Context, status, sort *string) ([]entities.Task, error) {
	var tasks []taskGORM
	db := r.db.BeginFind(ctx, &tasks)
	if status != nil {
		db = db.Where("status = ?", *status)
	}
	if sort != nil {
		db = db.OrderBy(*sort)
	}
	if err := db.Find(&tasks); err != nil {
		return nil, entities.ErrTaskFind
	}
	return r.convertToEntities(tasks), nil
}

func (r TaskRepo) DeleteTask(ctx context.Context, id uint32) error {
	if err := r.db.Delete(ctx, &taskGORM{}, &taskGORM{Id: id}); err != nil {
		return entities.ErrorTaskDelete
	}
	return nil
}

func (r TaskRepo) GetTask(ctx context.Context, id uint32) (entities.Task, error) {
	var task taskGORM
	if err := r.db.BeginFind(ctx, &taskGORM{}).Where("id = ?", id).Find(&task); err != nil {
		return entities.Task{}, entities.ErrTaskFind
	}
	return entities.NewTask(
		task.Id,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	), nil
}

func (r TaskRepo) convertToEntities(tasks []taskGORM) []entities.Task {
	var entitiesTasks []entities.Task
	for _, t := range tasks {
		entitiesTasks = append(entitiesTasks, entities.NewTask(
			t.Id,
			t.Title,
			t.Description,
			t.Status,
			t.CreatedAt,
			t.UpdatedAt,
		))
	}
	return entitiesTasks
}
