package utils

const STRIPE_STAGE_SK = ""
const STRIPE_PRODUCTION_SK = ""

const LOCAL = "local"
const STAGE = "stage"
const PROD = "production"

var STRIPE_CONST_LOOKUPS = map[string]map[string]string{
	LOCAL: {
		"REDIRECT_URL":   "http://localhost:3000/",
		"SECRET_KEY":     STRIPE_STAGE_SK,
		"PRO_PRODUCT_ID": "",
		"PRO_PRICE":      "",
	},
	STAGE: {
		"REDIRECT_URL":   "https://bip-app.vercel.app/",
		"SECRET_KEY":     STRIPE_STAGE_SK,
		"PRO_PRODUCT_ID": "",
		"PRO_PRICE":      "",
	},
	PROD: {
		"REDIRECT_URL":   "https://bip.so/",
		"SECRET_KEY":     STRIPE_PRODUCTION_SK,
		"PRO_PRODUCT_ID": "",
		"PRO_PRICE":      "",
	},
}
