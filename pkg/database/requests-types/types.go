package dbreq

type DBRequestType string

func (req DBRequestType) ToString() string {
	return string(req)
}

type DBRequest struct {
	Type DBRequestType
	Args []any
}

func (db_req DBRequest) ArgsIsValid(args_count int) bool {
	return len(db_req.Args) == args_count
}

func NewDBRequest(req_type DBRequestType, args ...any) DBRequest {
	return DBRequest{
		Type: req_type,
		Args: args,
	}
}
