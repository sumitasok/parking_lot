package models

type dataStore interface {
	Park(string, string) (*StoreyResponse, error)
	LeaveByPosition(int) (*StoreyResponse, error)
	FindByRegistrationNumber(string) (*StoreyResponse, error)
	FindAllByColor(string) (*StoreyResponse, error)
}