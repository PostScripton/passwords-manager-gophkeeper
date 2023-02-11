package models

type CredsSecret struct {
	ID             int64
	Website        string
	Login          string
	Password       string
	AdditionalData string
	UserID         int
}
