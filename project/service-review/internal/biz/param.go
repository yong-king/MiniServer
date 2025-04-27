package biz

type AuditParma struct {
	ReviewID int64
	Status   int32
	OpUser   string
	OpReason string
	OpMarks  string
}

type AppealParam struct {
	ReviewID  int64
	StoreID   int64
	Reason    string
	Content   string
	Picinfo   string
	VideoInfo string
	AppealID  int64
}

type ReplyParam struct {
	ReviewID  int64
	StoreID   int64
	Content   string
	Picinfo   string
	VideoInfo string
}

type AuditAppealParam struct{
	AppealID      int64               
    ReviewID      int64               
    Status        int32               
    OpUser        string  
}