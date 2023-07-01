package apiClient

import (
	"fmt"
	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/pkg/configs"
	"log"
	"time"
)

// AsyncScheduler for ref: https://github.com/hibiken/asynq/wiki/Periodic-Tasks
func AsyncScheduler() {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("error in reading time location", err)
	}
	redisAddr := fmt.Sprintf("%s:%s", configs.GetRedisConfig().Host, configs.GetRedisConfig().Port)
	scheduler := asynq.NewScheduler(asynq.RedisClientOpt{Addr: redisAddr, Password: configs.GetRedisConfig().Password}, &asynq.SchedulerOpts{
		Location: loc,
	})

	// Initialize the tasks here.
	task1 := asynq.NewTask(TestTaskCronMethod1, nil)
	canvasBranchAccessCron := asynq.NewTask(CanvasBranchAccessCron, nil)
	runFailedDiscordEventsCron := asynq.NewTask(RunFailedDiscordEventsCron, nil)

	// Register the task.
	// You can use "@every <duration>" to specify the interval.
	entryID1, err := scheduler.Register("20 15 * * *", task1)
	canvasBranchAccessID, err := scheduler.Register("0 9 * * *", canvasBranchAccessCron)
	_, err = scheduler.Register("00 01,13 * * *", runFailedDiscordEventsCron)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("registered an entry: %q\n%q\n", entryID1, canvasBranchAccessID)

	// Starting the scheduler.
	if err = scheduler.Run(); err != nil {
		log.Fatal(err)
	}
}
