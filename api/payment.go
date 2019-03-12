package api

import (
	"time"

	"github.com/google/uuid"
)

type PaymentParty struct {
	AccountName       string `json:"accountName"`
	AccountNumber     string `json:"accountNumber"`
	AccountNumberCode string `json:"accountNumberCode"`
	BankID            string `json:"bankId"`
	BankIDCode        string `json:"bankIdCode"`
	Name              string `json:"name"`
	Address           string `json:"address"`
}

type Payment struct {
	ID                   uuid.UUID     `json:"paymentId"`
	CreatedAt            *time.Time    `json:"createdAt"`
	UpdatedAt            *time.Time    `json:"updatedAt"`
	Purpose              string        `json:"purpose"`
	Scheme               string        `json:"scheme"`
	Type                 string        `json:"type"`
	Amount               string        `json:"amount"`
	Beneficiary          *PaymentParty `json:"beneficiary"`
	Currency             string        `json:"currency"`
	DebitorParty         *PaymentParty `json:"debitorParty"`
	EndToEndReference    string        `json:"endToEndReference"`
	NumericReference     string        `json:"numericReference"`
	ProcessingDate       string        `json:"processingDate"`
	Reference            string        `json:"reference"`
	SchemePaymentSubType string        `json:"schemePaymentSubType"`
	SchemePaymentType    string        `json:"schemePaymentType"`
	ChargesInformation   struct {
		BearerCode              string `json:"bearerCode"`
		ReceiverChargesAmount   string `json:"receiverChargesAmount"`
		ReceiverChargesCurrency string `json:"receiverChargesCurrency"`
		SenderCharges           []struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency"`
		} `json:"senderCharges"`
	} `json:"chargesInformation"`
	FX struct {
		ContractReference string `json:"contractReference"`
		ExchangeRate      string `json:"exchangeRate"`
		OriginalAmount    string `json:"originalAmount"`
		OriginalCurrency  string `json:"originalCurrency"`
	} `json:"fx"`
}

func NewPayment() *Payment {
	now := time.Now()
	return &Payment{
		ID:        uuid.New(),
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}

func (p *Payment) SetScheme(value string) *Payment {
	p.Scheme = value
	return p
}

func (p *Payment) Validate() error {
	return nil
}
