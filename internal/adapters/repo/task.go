package repo

import (
	"context"
	"test-ns/internal/entities"
	"time"

	"gorm.io/gorm/clause"
)

type taskGORM struct {
	ID          uint32 `gorm:"primaryKey;column:id"`
	Title       string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   *time.Time `gorm:"autoUpdateTime:false"`
}

func (taskGORM) TableName() string {
	return "task_gorms"
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
		ID:          task.ID(),
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
		t.ID,
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
		"updated_at":  task.UpdatedAt(),
	}
	if err := r.db.Update(ctx, "task_gorms", &updates, &taskGORM{ID: task.ID()}); err != nil {
		return entities.Task{}, entities.ErrorTaskUpdate
	}
	return entities.NewTask(
		task.ID(),
		updates["title"].(string),
		updates["description"].(string),
		updates["status"].(string),
		task.CreatedAt(),
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
	if err := r.db.Delete(ctx, &taskGORM{}, &taskGORM{ID: id}); err != nil {
		return entities.ErrorTaskDelete
	}
	return nil
}

func (r TaskRepo) GetTask(ctx context.Context, id uint32) (entities.Task, error) {
	var task taskGORM
	if err := r.db.BeginFind(ctx, &taskGORM{}).Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).First(&task); err != nil {
		return entities.Task{}, entities.ErrTaskFind
	}
	return entities.NewTask(
		task.ID,
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
			t.ID,
			t.Title,
			t.Description,
			t.Status,
			t.CreatedAt,
			t.UpdatedAt,
		))
	}
	return entitiesTasks
}
