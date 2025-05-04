package memory

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"pinstack-user-service/internal/custom_errors"
	"pinstack-user-service/internal/logger"
	"pinstack-user-service/internal/model"
)

type Repository struct {
	users  map[int64]*model.User
	mu     sync.RWMutex
	nextID int64
	log    *logger.Logger
}

func NewUserRepository(log *logger.Logger) *Repository {
	return &Repository{
		users:  make(map[int64]*model.User),
		nextID: 1,
		log:    log,
	}
}

func (r *Repository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Username == user.Username {
			return nil, custom_errors.ErrUsernameExists
		}
		if u.Email == user.Email {
			return nil, custom_errors.ErrEmailExists
		}
	}

	now := time.Now()
	user.ID = r.nextID
	user.CreatedAt = now
	user.UpdatedAt = now
	r.users[user.ID] = user
	r.nextID++

	return user, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, custom_errors.ErrUserNotFound
	}

	return user, nil
}

func (r *Repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}

	return nil, custom_errors.ErrUserNotFound
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, custom_errors.ErrUserNotFound
}

func (r *Repository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existingUser, exists := r.users[user.ID]
	if !exists {
		return nil, custom_errors.ErrUserNotFound
	}

	for _, u := range r.users {
		if u.ID != user.ID {
			if u.Username == user.Username && u.Username != existingUser.Username {
				return nil, custom_errors.ErrUsernameExists
			}
			if u.Email == user.Email && u.Email != existingUser.Email {
				return nil, custom_errors.ErrEmailExists
			}
		}
	}

	user.CreatedAt = existingUser.CreatedAt
	user.UpdatedAt = time.Now()
	r.users[user.ID] = user

	return user, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return custom_errors.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func (r *Repository) Search(ctx context.Context, searchQuery string, offset, pageSize int) ([]*model.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []*model.User
	for _, user := range r.users {
		if searchQuery == "" ||
			contains(user.Username, searchQuery) ||
			contains(user.Email, searchQuery) ||
			(user.FullName != nil && contains(*user.FullName, searchQuery)) ||
			(user.Bio != nil && contains(*user.Bio, searchQuery)) {
			results = append(results, user)
		}
	}

	// Сортируем результаты по username для стабильной пагинации
	sort.Slice(results, func(i, j int) bool {
		return results[i].Username < results[j].Username
	})

	total := len(results)
	if offset >= total {
		return []*model.User{}, total, nil
	}

	end := offset + pageSize
	if end > total {
		end = total
	}

	return results[offset:end], total, nil
}

func (r *Repository) UpdatePassword(ctx context.Context, id int64, password string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return custom_errors.ErrUserNotFound
	}

	user.Password = password
	user.UpdatedAt = time.Now()
	return nil
}

func (r *Repository) UpdateAvatar(ctx context.Context, id int64, avatarURL string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return custom_errors.ErrUserNotFound
	}

	user.AvatarURL = &avatarURL
	user.UpdatedAt = time.Now()
	return nil
}

func contains(s, substr string) bool {
	return len(substr) == 0 || strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}
