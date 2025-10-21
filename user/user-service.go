package user

type Service struct {
	*baseService
}

type baseService struct {
}

func NewService() *Service {
	return &Service{
	}
}

func (b *baseService) GetUsers() ([]User, error) {
	return getUsers()
}

func (b *baseService) GetUserByEmail(email string) (*User, error) {
	return getUserByEmail(email)
}

func (b *baseService) GetRoleByUID(uid string) (*Role, error) {
	return getRoleByUID(uid)
}
