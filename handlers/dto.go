package handlers

type createAgreementResp struct {
	AgreementID string `json:"agreement_id"`
}

type statusAgreementResp struct {
	Success bool `json:"success"`
}

type createApartmentResp struct {
	ApartmentID string `json:"apartment_id"`
}

type statusApartmentResp struct {
	Success bool `json:"success"`
}
