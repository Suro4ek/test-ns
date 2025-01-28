package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-ns/internal/adapters/api/action"
	"test-ns/internal/adapters/presenter"
	"test-ns/internal/adapters/repo"
	"test-ns/internal/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

const timeout = time.Second * 5

type Router struct {
	transaction repo.GTransaction
	db          repo.GSQL
	gRouter     *gin.Engine
}

func NewRouter(transaction repo.GTransaction, db repo.GSQL) Router {
	return Router{
		transaction: transaction,
		db:          db,
		gRouter:     gin.Default(),
	}
}

func (r Router) Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.Recovery()
	r.register()

	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		Addr:         ":8080",
		Handler:      r.gRouter,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			panic(err)
		}
	}()

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
}

func (r Router) register() {
	r.gRouter.POST("/tasks", r.createTaskAction())
	r.gRouter.GET("/tasks", r.getTasksAction())
	r.gRouter.PUT("/tasks/:id", r.updateTaskAction())
	r.gRouter.DELETE("/tasks/:id", r.deleteTaskAction())
}

func (r Router) createTaskAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			uc = usecase.NewCreateTaskInteractor(
				repo.NewTaskRepo(r.db),
				presenter.NewCreateTaskPresenter(),
				timeout,
			)

			act = action.NewCreateTaskAction(uc)
		)

		act.Execute(ctx.Writer, ctx.Request)
	}
}

func (r Router) getTasksAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			uc = usecase.NewTaskInteractor(
				repo.NewTaskRepo(r.db),
				presenter.NewGetTasksPresenter(),
				timeout,
			)

			act = action.NewGetTasksAction(uc)
		)

		act.Execute(ctx.Writer, ctx.Request)
	}
}

func (r Router) updateTaskAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			uc = usecase.NewUpdateTaskInteractor(
				r.transaction,
				repo.NewTaskRepo(r.db),
				presenter.NewUpdateTaskPresenter(),
				timeout,
			)

			act = action.NewUpdateTaskAction(uc)
		)
		q := ctx.Request.URL.Query()
		q.Add("id", ctx.Param("id"))
		ctx.Request.URL.RawQuery = q.Encode()
		act.Execute(ctx.Writer, ctx.Request)
	}
}

func (r Router) deleteTaskAction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			uc = usecase.NewDeleteTaskInteractor(
				repo.NewTaskRepo(r.db),
				timeout,
			)

			act = action.NewDeleteTaskAction(uc)
		)
		q := ctx.Request.URL.Query()
		q.Add("id", ctx.Param("id"))
		ctx.Request.URL.RawQuery = q.Encode()
		act.Execute(ctx.Writer, ctx.Request)
	}
}
