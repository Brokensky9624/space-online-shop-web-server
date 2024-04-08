package member

type MemberParam struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type MemberEditParam struct {
	Username string `json:"username"`
	// Password    string `json:"password"`
	// NewPassword string `json:"newpassword"`
	Email string `json:"email"`
}

type MemberDeleteParam struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type MemberInfoParam struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Birth    string `json:"birth"`
	Phone    string `json:"phone"`
}
