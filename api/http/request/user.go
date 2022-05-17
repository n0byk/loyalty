package request

type UserRegistration struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLogin struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserWithdraw struct {
	Order string  `json:"order"`
	Sum   float32 `json:"sum"`
}
