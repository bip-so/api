package models

const FE_LOCAL = "local"
const FE_STAGE = "stage"
const FE_PROD = "production"

var MailerRouterPaths = map[string]map[string]string{
	FE_LOCAL: {
		"BASE_URL": "http://localhost:3000/",
	},
	FE_STAGE: {
		"BASE_URL": "https://bip-app.vercel.app/",
	},
	FE_PROD: {
		"BASE_URL": "https://bip.so/",
	},
}
