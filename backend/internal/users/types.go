package users

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Role      Role      `json:"role"`    // e.g., "admin", "user"
	Active    bool      `json:"active"`  // Indicates if the user account is active
	Profile   Profile   `json:"profile"` // Additional user profile information
}

type Profile struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Bio       string `json:"bio"`        // Short biography or description
	AvatarURL string `json:"avatar_url"` // URL to the user's avatar image
	Links     []Link `json:"links"`      // List of social media or other relevant links
}

type Link struct {
	URL    string `json:"url"`    // URL of the link
	Title  string `json:"title"`  // Title or description of the link
	Icon   string `json:"icon"`   // Optional icon or image URL for the link
	Active bool   `json:"active"` // Indicates if the link is currently active
}

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
	RoleGuest Role = "guest"
)
