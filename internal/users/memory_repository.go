package users

import (
	"errors"
)

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

func (repo *memoryRepository) Save(user *user, identity *Identity) error {
	if user == nil {
		return errors.New("user is nil")
	}

	key := identity.Provider + ":" + identity.ExternalID
	if _, exists := repo.identities[key]; exists {
		return errors.New("identity has been taken")
	}

	repo.users[user.ID] = user
	repo.identities[key] = user.ID

	return nil
}

func (repo *memoryRepository) AddIdentity(userID string, identity *Identity) error {
	if _, exists := repo.users[userID]; !exists {
		return errors.New("user not found")
	}

	key := identity.Provider + ":" + identity.ExternalID
	if _, exists := repo.identities[key]; exists {
		return errors.New("identity has been taken")
	}

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

func (repo *memoryRepository) GetByExternalID(identity *Identity) (*user, error) {
	key := identity.Provider + ":" + identity.ExternalID
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

func (repo *memoryRepository) getUserCount() (int, error) {
	return len(repo.users), nil
}

var _ UserRepository = (*memoryRepository)(nil)
