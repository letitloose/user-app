package user

type UserService struct {
	repository *userRepository
}

func NewUserService(repository *userRepository) *UserService {
	return &UserService{repository: repository}
}

func (service *UserService) ListAllUsers() []*User {
	return service.repository.listAll()
}

func (service *UserService) FindByUsername(username string) (*User, error) {
	return service.repository.findUser(username)
}

func (service *UserService) RemoveUser(username string) error {
	return service.repository.removeUser(username)
}

func (service *UserService) AddUser(user *User) error {
	return service.repository.addUser(user)
}

func (service *UserService) UpdateUser(user *User) error {
	return service.repository.updateUser(user)
}
