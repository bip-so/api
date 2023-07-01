package tasks

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"gitlab.com/phonepost/bip-be-platform/internal/models"
	"gitlab.com/phonepost/bip-be-platform/pkg"
)

func (s taskService) HandleLoginEmailTask(ctx context.Context, task *asynq.Task) error {
	var user models.User
	json.Unmarshal(task.Payload(), &user)
	body := "<div><strong> Your account with username " + user.Username + " has logged in at " + user.UpdatedAt.String() + "</strong></div>"
	bodyPlainText := "<div><strong> Your account with username " + user.Username + " has logged in at " + user.UpdatedAt.String() + "</strong></div>"
	subject := "LoggedIn confirmation email"
	var mailer pkg.BipMailer
	toList := []string{user.Email.String}
	emptyList := []string{}
	err := pkg.BipMailer.SendEmail(mailer, toList, emptyList, emptyList, subject, body, bodyPlainText)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
