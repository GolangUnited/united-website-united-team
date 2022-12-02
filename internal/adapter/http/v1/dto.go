package v1

type PostUserLoginOutput struct {
	UserId      string `json:"user_id"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}
