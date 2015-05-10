package oneandone_cloudserver_api

import ()

type User struct {
	withId
	withName
	State     string `json:"state"`
	Role      string `json:"role"`
	ApiActive bool   `json:"api_active"`
}

type UserCreateData struct {
	withName
	Password string `json:"password"`
	Email    string `json:"email"`
}

// GET /users
func (api *API) GetUsers() []User {
	return []User{}
}

// POST /users
func (api *API) CreateUser(configuration UserCreateData) User {
	return User{}
}

// GET /users/{id}
func (api *API) GetUser(Id string) User {
	return User{}
}

// DELETE /users/{id}
func (user *User) Delete() User {
	return User{}
}

// PUT /users/{id}
func (user *User) ModifyUser(data UserCreateData) User {
	return User{}
}

// GET /users/{id}/api

// PUT /users/{id}/api

// GET /users/{id}/api/ips

// POST /users/{id}/api/ips

// GET /users/{id}/api/ips/{ipId}

// DELETE /users/{id}/api/ips/{ipId}


