package email

import (
	"github.com/matcornic/hermes/v2"
)

type confirm struct {
}

func (r *confirm) Name() string {
	return "confirm"
}

func (r *confirm) Email(name, link string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"Tap the button below to confirm your email address.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Color: "#DC4D2F",
						Text:  "Confirm your account",
						Link:  link,
					},
				},
			},
			Outros: []string{
				"If you didn't create an account with Saferwall, you can safely delete this email..",
			},
			Signature: "Thanks",
		},
	}
}
