package tasks

import (
	"context"
	"fmt"
	"github.com/hibiken/asynq"
	"time"
)

func (s taskService) TestTaskCronMethod(ctx context.Context, task *asynq.Task) error {
	fmt.Println("Printing the time TestTaskCronMethod =========>", time.Now())
	return nil
}

func (s taskService) TestTaskCronMethod1(ctx context.Context, task *asynq.Task) error {
	fmt.Println("Printing the time TestTaskCronMethod1 =========>", time.Now())
	return nil
}
