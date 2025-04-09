package models

import (
	"fmt"
	"strings"
)

type Payment struct {
	OrderID   string
	Currency  Currency
	Amount    int64
	Signature string
	Payload   Payload
}

type Payload map[string]string

func (p Payload) Generate() string {
	if p == nil {
		return ""
	}

	builder := &strings.Builder{}
	builder.WriteString("&")

	for key, value := range p {
		param := fmt.Sprintf("us_%s=%s&", key, value)
		builder.WriteString(param)
	}

	return builder.String()[:builder.Len()-1]
}
