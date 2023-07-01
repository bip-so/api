package mailers

type SendEmail struct {
	userIDs   []uint64
	ToEmails  []string
	CcEmails  []string
	BccEmails []string
	Subject   string
	BodyHtml  string
	BodyText  string
	ExtraData string
}
