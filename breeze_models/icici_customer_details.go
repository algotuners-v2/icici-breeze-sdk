package breeze_models

type CustomerDetailsResponse struct {
	Success CustomerDetailsBody `json:"Success"`
	Status  int                 `json:"Status"`
	Error   any                 `json:"Error"`
}

type CustomerDetailsBody struct {
	IdirectUserid   string `json:"idirect_userid"`
	SessionToken    string `json:"session_token"`
	IdirectUserName string `json:"idirect_user_name"`
}
