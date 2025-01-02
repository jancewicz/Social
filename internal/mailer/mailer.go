package mailer

import "embed"

const (
	FromName           = "GopherSocial"
	maxRetires         = 3
	InvitationTemplate = "user_invitation.tmpl"
)

//go:embed "templates"
var FS embed.FS

type Client interface {
	// isSandbox checks if environment is dev or prod to avoid sending emails from dev
	Send(template, username, email string, data any, isSendbox bool) error
}
