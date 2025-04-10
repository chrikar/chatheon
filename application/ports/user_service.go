package ports

//go:generate mockery --name=UserService --output=../../adapters/mocks --case=underscore
type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}
