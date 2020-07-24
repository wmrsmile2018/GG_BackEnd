package model

import "testing"

func TestUser(t *testing.T) *User {
	return &User{
		ID:       "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a13",
		Email:    "test@test.org",
		Password: "password",
	}
}
