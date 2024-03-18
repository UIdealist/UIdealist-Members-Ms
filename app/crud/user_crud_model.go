package crud

// UserCreate struct to describe create user object.
type UserCreate struct {
	Username string `json:"username" validate:"required,username,lte=255"`
	Email    string `json:"email" validate:"required,lte=255"`
}

// AnonymousUserCreate struct to describe create anonymous user object.
type AnonymousUserCreate struct {
	TempName string `json:"temp_name" validate:"required,lte=255"`
}
