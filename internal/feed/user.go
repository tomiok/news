package feed

type UserStorage interface {
	SaveOrValidate(emailOrToken string) error
}

type UserService struct {
	UserStorage
}

func (u *UserService) UserLogin(emailOrToken string) error {
	return u.SaveOrValidate(emailOrToken)
}
