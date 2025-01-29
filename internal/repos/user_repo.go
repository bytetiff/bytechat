package repositories

import "github.com/bytetiff/bytechat/internal/models"

type UserRepository struct {
	storage map[int]models.User
	lastID  int
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		storage: make(map[int]models.User),
		lastID:  0,
	}
}

func (r *UserRepository) CreateUser(u models.User) (models.User, error) {
	r.lastID++
	u.ID = r.lastID
	r.storage[u.ID] = u
	return u, nil
}

func (r *UserRepository) GetUserByUsername(username string) (models.User, bool) {
	for _, user := range r.storage {
		if user.Username == username {
			return user, true
		}
	}
	return models.User{}, false
}
