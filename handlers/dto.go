package handlers

type createAgreementResp struct {
	AgreementID string `json:"agreement_id"`
}

type statusAgreementReq struct {
	Action string `json:"action"`
}

type statusAgreementResp struct {
	Success bool `json:"success"`
}

type createApartmentResp struct {
	ApartmentID string `json:"apartment_id"`
}

type statusApartmentReq struct {
	Action string `json:"action"`
}

type statusApartmentResp struct {
	Success bool `json:"success"`
}
