package notifications

import "fmt"

func (s notificationService) GenerateEmailTemplate(text string, url string, event string) string {
	var template string
	if event == "CanvasLimitExceed" {
		template = `<div><ul>
			<li>You have exceeded the 25 private canvas limit of the free plan. Please upgrade to create more private canvases!</li>
			<li>Newly published canvases will be made public as you have exhausted your free plan limits. Upgrade now to keep them private!</li>
			</ul>` + "<br><br>" + fmt.Sprintf("<a href=%s>Upgrade</a></div>", url)
	} else {
		template = text + "<br><br>" + fmt.Sprintf("<a href=%s>click here</a>", url)
	}
	return template
}
