package tool

import "golang.org/x/crypto/bcrypt"

func HashPassword(pwd string) (string, error) {
	ret, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

func CheckPassword(uncheckPwd, hashedPwd string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(uncheckPwd))
	return err == nil
}
