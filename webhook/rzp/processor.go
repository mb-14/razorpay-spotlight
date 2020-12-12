package rzp

import (
	"strings"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"github.com/mb-14/rzp-spotlight/webhook/json"
)

// Events
const (
	PaymentAuthorized = "payment.authorized"
	PaymentFailed     = "payment.failed"
)

// Methods
const (
	Netbanking = "netbanking"
	Wallet     = "wallet"
	UPI        = "upi"
	Card       = "card"
)

func ProcessPayloadJson(json json.Json) *write.Point {
	var event string
	if value, err := json.GetString("event"); err == nil {
		event = value
	}
	if event == PaymentAuthorized {
		return generatePaymentEvent(json, "payment_authorized")
	} else if event == PaymentFailed {
		return generatePaymentEvent(json, "payment_failed")
	}
	return nil
}

func generatePaymentEvent(json json.Json, event string) *write.Point {
	amount, _ := json.GetInt("payload.payment.entity.amount")
	createdAt, _ := json.GetTime("payload.payment.entity.created_at")
	return influxdb2.NewPoint(event,
		addTags(json, event),
		map[string]interface{}{"amount": amount},
		createdAt)
}

func addTags(p json.Json, event string) map[string]string {
	tags := make(map[string]string)
	// Common tags
	method, _ := p.GetString("payload.payment.entity.method")
	currency, _ := p.GetString("payload.payment.entity.currency")
	tags["currency"] = currency
	tags["method"] = method
	if method == Netbanking {
		tags["bank"], _ = p.GetString("payload.payment.entity.bank")
	}

	if method == Wallet {
		tags["walletName"], _ = p.GetString("payload.payment.entity.wallet")
	}

	if method == UPI {
		vpa, _ := p.GetString("payload.payment.entity.vpa")
		vpaString := strings.Split(vpa, "@")
		tags["upiPsp"] = vpaString[1]
	}

	if method == Card {
		tags["cardNetwork"], _ = p.GetString("payload.payment.entity.card.network")
		tags["cardType"], _ = p.GetString("payload.payment.entity.card.type")
		tags["cardInternational"], _ = p.GetString("payload.payment.entity.card.international")
		tags["cardIssuer"], _ = p.GetString("payload.payment.entity.card.issuer")
	}

	if event == "payment_failed" {
		tags["errorSource"], _ = p.GetString("payload.payment.entity.error_source")
	}

	return tags
}
