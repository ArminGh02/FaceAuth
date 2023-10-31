package model

type RegisterStatus int

const (
	StatusPending RegisterStatus = iota
	StatusRejected
	StatusAccepted
)

func (rs RegisterStatus) String() string {
	switch rs {
	case StatusPending:
		return "pending"
	case StatusRejected:
		return "rejected"
	case StatusAccepted:
		return "accepted"
	default:
		return "unknown"
	}
}
