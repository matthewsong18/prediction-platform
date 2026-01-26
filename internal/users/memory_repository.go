package users

import "errors"

type memoryRepository struct {
	users      map[string]*user
	identities map[string]string // Key: provider:externalID, Value: userID
}

func NewMemoryRepository() UserRepository {
	return &memoryRepository{
		users:      make(map[string]*user),
		identities: make(map[string]string),
	}
}

func (repo *memoryRepository) Save(user *user) error {
	if user == nil {
		return errors.New("user is nil")
	}

	repo.users[user.ID] = user
	return nil
}

func (repo *memoryRepository) AddIdentity(userID, provider, externalID string) error {
	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}
	key := provider + ":" + externalID
	repo.identities[key] = userID
	return nil
}

func (repo *memoryRepository) GetByID(id string) (*user, error) {
	user, exists := repo.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (repo *memoryRepository) GetByExternalID(provider, externalID string) (*user, error) {
	key := provider + ":" + externalID
	userID, exists := repo.identities[key]
	if !exists {
		return nil, errors.New("user not found")
	}
	return repo.GetByID(userID)
}

func (repo *memoryRepository) Delete(userID string) error {
	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	delete(repo.users, userID)

	// Clean up identities (inefficient but fine for memory repo)
	for k, v := range repo.identities {
		if v == userID {
			delete(repo.identities, k)
		}
	}
	return nil
}

var _ UserRepository = (*memoryRepository)(nil)
