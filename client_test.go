package freekassa

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenorachi/freekassa-sdk-go/models"
)

func TestClient_GenerateInvoice(t *testing.T) {
	t.Parallel()

	type args struct {
		merchantID int64
		p          *models.Payment
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				merchantID: 1,
				p: &models.Payment{
					OrderID:   "order1",
					Currency:  models.USD,
					Amount:    10,
					Signature: "sign1234",
					Payload: map[string]string{
						"key1": "value1",
					},
				},
			},
			want: "https://pay.fk.money/?m=1&o=order1&oa=10&currency=USD&s=sign1234&us_key1=value1",
		},
		{
			name: "OK_NilPayment",
			args: args{
				merchantID: 1,
				p:          nil,
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewClient(&Config{
				MerchantID: tt.args.merchantID,
				SecretKey1: "",
				SecretKey2: "",
			})

			assert.Equal(t, tt.want, c.GenerateInvoice(tt.args.p))
		})
	}
}

func TestClient_GenerateInvoiceSignature(t *testing.T) {
	t.Parallel()

	type args struct {
		amount     int64
		currency   models.Currency
		orderID    string
		merchantID int64
		secretKey1 string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				amount:     10,
				currency:   models.EUR,
				orderID:    "order1",
				merchantID: 199,
				secretKey1: "secretkey1",
			},
			want: md5Hash("199:10:secretkey1:EUR:order1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewClient(&Config{
				MerchantID: tt.args.merchantID,
				SecretKey1: tt.args.secretKey1,
			})

			assert.Equal(t, tt.want, c.GenerateInvoiceSignature(tt.args.amount, tt.args.currency, tt.args.orderID))
		})
	}
}

func TestClient_GenerateConfirmSignature(t *testing.T) {
	t.Parallel()

	type args struct {
		amount     int64
		orderID    string
		merchantID int64
		secretKey2 string
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "OK",
			args: args{
				amount:     1337,
				orderID:    "order1",
				merchantID: 777,
				secretKey2: "secretkey2",
			},
			want: md5Hash("777:1337:secretkey2:order1"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewClient(&Config{
				MerchantID: tt.args.merchantID,
				SecretKey2: tt.args.secretKey2,
			})

			assert.Equal(t, tt.want, c.GenerateConfirmSignature(tt.args.amount, tt.args.orderID))
		})
	}
}
