package ports

type UserService interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}
