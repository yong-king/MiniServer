package biz

type AuditReviewParam struct {
	ReviewID  int64
	Status    int32
	OpUser    string
	OpReason  string
	OpRemarks *string
}

type AuditAppealParam struct {
	AppealID  int64
	ReviewID  int64
	Status    int32
	OpUser    string
	OpReason  string
	OpRemarks *string
}
