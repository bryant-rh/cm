package client

type UserData struct {
	GlobalRes
	Data string `json:"data"`
}

//User_Login /user/login
func (c *CMClient) User_Login(username, password string) (*UserData, error) {
	res := &UserData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"username": username,
			"password": password,
		}).
		SetResult(res).Get("/user/login")
	return res, err
}

//User_VerifyToken /user/verifyToken
func (c *CMClient) User_VerifyToken(token string) (*UserData, error) {
	res := &UserData{}
	_, err := c.R().
		SetQueryParams(map[string]string{ // Set multiple query params at once
			"token": token,
		}).
		SetResult(res).Get("/user/verifytoken")
	return res, err
}
