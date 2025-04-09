package freekassa

import (
	"fmt"
)

type IClient interface {
	GenerateInvoice(p *Payment) string
	GenerateInvoiceSignature(amount int64, currency Currency, orderID string) string
	GenerateConfirmSignature(amount int64, orderID string) string
}

type Config struct {
	MerchantID int64
	SecretKey1 string
	SecretKey2 string
}

type client struct {
	cfg *Config
}

func NewClient(cfg *Config) IClient {
	return &client{cfg: cfg}
}

func (c *client) GenerateInvoice(p *Payment) string {
	if p == nil {
		return ""
	}

	return fmt.Sprintf(
		"%s/?m=%d&o=%s&oa=%d&currency=%s&s=%s%s",
		InvoiceBaseURL,
		c.cfg.MerchantID,
		p.OrderID,
		p.Amount,
		p.Currency.String(),
		p.Signature,
		p.Payload.Generate(),
	)
}

func (c *client) GenerateInvoiceSignature(amount int64, currency Currency, orderID string) string {
	signData := fmt.Sprintf("%d:%d:%s:%s:%s", c.cfg.MerchantID, amount, c.cfg.SecretKey1, currency.String(), orderID)

	return md5Hash(signData)
}

func (c *client) GenerateConfirmSignature(amount int64, orderID string) string {
	signData := fmt.Sprintf("%d:%d:%s:%s", c.cfg.MerchantID, amount, c.cfg.SecretKey2, orderID)

	return md5Hash(signData)
}
