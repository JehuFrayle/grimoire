package users

import (
	"context"
	"fmt"
	"time"
)

type MemUserRepository struct {
	users map[string]User // In-memory storage for users
}

var initialUsers = map[string]User{
	"u1": {
		ID:        "u1",
		Username:  "jdoe",
		Email:     "jdoe@example.com",
		CreatedAt: time.Now().Add(-48 * time.Hour),
		UpdatedAt: time.Now().Add(-2 * time.Hour),
		Role:      RoleAdmin,
		Active:    true,
		Profile: Profile{
			FirstName: "John",
			LastName:  "Doe",
			Bio:       "System administrator and DevOps enthusiast.",
			AvatarURL: "https://example.com/avatars/jdoe.png",
			Links: []Link{
				{URL: "https://github.com/jdoe", Title: "GitHub", Icon: "github", Active: true},
				{URL: "https://linkedin.com/in/jdoe", Title: "LinkedIn", Icon: "linkedin", Active: true},
			},
		},
	},
	"u2": {
		ID:        "u2",
		Username:  "asmith",
		Email:     "asmith@example.com",
		CreatedAt: time.Now().Add(-72 * time.Hour),
		UpdatedAt: time.Now().Add(-24 * time.Hour),
		Role:      RoleUser,
		Active:    true,
		Profile: Profile{
			FirstName: "Alice",
			LastName:  "Smith",
			Bio:       "Frontend developer and UX designer.",
			AvatarURL: "https://example.com/avatars/asmith.jpg",
			Links: []Link{
				{URL: "https://twitter.com/asmith", Title: "Twitter", Icon: "twitter", Active: true},
			},
		},
	},
	"u3": {
		ID:        "u3",
		Username:  "guest1",
		Email:     "guest1@example.com",
		CreatedAt: time.Now().Add(-10 * time.Hour),
		UpdatedAt: time.Now().Add(-1 * time.Hour),
		Role:      RoleGuest,
		Active:    false,
		Profile: Profile{
			FirstName: "Guest",
			LastName:  "User",
			Bio:       "Limited access guest account.",
			AvatarURL: "",
			Links:     []Link{},
		},
	},
}

func NewMemUserRepository() *MemUserRepository {
	return &MemUserRepository{
		users: initialUsers,
	}
}
func (r *MemUserRepository) GetAll(ctx context.Context) ([]User, error) {
	var userList []User
	for _, user := range r.users {
		userList = append(userList, user)
	}
	return userList, nil
}
func (r *MemUserRepository) GetByID(ctx context.Context, id string) (*User, error) {
	user, exists := r.users[id]
	if !exists {
		return nil, fmt.Errorf("user with id %s not found", id)
	}
	return &user, nil
}
func (r *MemUserRepository) Create(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; exists {
		return fmt.Errorf("user with id %s already exists", user.ID)
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	if user.ID == "" {
		user.ID = fmt.Sprintf("u%d", len(r.users)+1) // Simple ID generation
	}
	if user.Role == "" {
		user.Role = RoleUser // Default role if not specified
	}

	r.users[user.ID] = *user
	return nil
}
func (r *MemUserRepository) Update(ctx context.Context, user *User) error {
	if _, exists := r.users[user.ID]; !exists {
		return fmt.Errorf("user with id %s not found", user.ID)
	}
	r.users[user.ID] = *user
	return nil
}
func (r *MemUserRepository) Delete(ctx context.Context, id string) error {
	if _, exists := r.users[id]; !exists {
		return fmt.Errorf("user with id %s not found", id)
	}
	delete(r.users, id)
	return nil
}
