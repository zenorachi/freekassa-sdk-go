package models

const InvoiceBaseURL = `https://pay.fk.money`

func IPsWhitelist() map[string]struct{} {
	return map[string]struct{}{
		"168.119.157.136": {},
		"168.119.60.227":  {},
		"178.154.197.79":  {},
		"51.250.54.238":   {},
	}
}
