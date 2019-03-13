package api

// A PaymentID parameter model.
//
// This is used for operations that want the ID of an pet in the path
// swagger:parameters getPayment deletePayment savePayment
type paymentID struct {
	// The ID of the payment
	//
	// in: path
	ID string `json:"id"`
}

// List of payments with paging info
// swagger:response paymentList
type paymentList struct {
	// in: body
	Body struct {
		Results  []Payment `json:"results"`
		Total    int       `json:"total"`
		SubTotal int       `json:"subTotal"`
	}
}

// swagger:response reqError
type reqError struct {
	// in: body
	Body struct {
		APIError
	}
}

// swagger:response singlePayment
type singlePayment struct {
	// in: body
	Body struct {
		Payment
	}
}
