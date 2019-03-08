package api

const DocumentTypePayment = "Payment"

type paymentParty struct {
	AccountName       string `json:"accountName"`
	AccountNumber     string `json:"accountNumber"`
	AccountNumberCode string `json:"accountNumberCode"`
	BankID            string `json:"bankId"`
	BankIDCode        string `json:"bankIdCode"`
	Name              string `json:"name"`
	Address           string `json:"address"`
}

type Payment struct {
	Amount               string        `json:"amount"`
	Beneficiary          *paymentParty `json:"beneficiary"`
	Currency             string        `json:"currency"`
	DebitorParty         *paymentParty `json:"debitorParty"`
	EndToEndReference    string        `json:"endToEndReference"`
	NumericReference     string        `json:"numericReference"`
	PaymentID            string        `json:"paymentId"`
	PaymentPurpose       string        `json:"paymentPurpose"`
	PaymentScheme        string        `json:"paymentScheme"`
	PaymentType          string        `json:"paymentType"`
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
