package request

type ClaimStudentRequest struct {
	Code     string `json:"code"`
	ParentId string `json:"parent_id"`
}
