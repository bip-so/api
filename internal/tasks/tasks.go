package tasks

import (
	"context"
	"fmt"
	"reflect"

	"github.com/hibiken/asynq"
)

func DefaultTaskHandler(ctx context.Context, task *asynq.Task) error {
	taskName := task.Type()[6:]
	fmt.Println("Received task", taskName)
	s := taskService{}
	method := reflect.ValueOf(s).MethodByName(taskName)
	params := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(task)}
	method.Call(params)
	return nil
}
