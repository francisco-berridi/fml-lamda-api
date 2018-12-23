package authentication

import "domain"

func CreateNewUser(email, firstName, lastName, password string) (domain.User, error) {

	// Validate password
	//password

	return domain.User{
		Id:        "aFgf98df3",
		Username:  email,
		FirstName: firstName,
		LastName:  lastName,
	}, nil
}
