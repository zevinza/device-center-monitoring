package lib

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

func saltAes(plain string) string {
	return plain + "//" + viper.GetString("SALT") + "//" + viper.GetString("AES")
}

// PasswordEncrypt Password Encrypt
func PasswordEncrypt(plain string) string {
	password := saltAes(plain)
	if hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); nil == err {
		return string(hashed)
	}

	return ""
}

// PasswordCompare Password Compare
func PasswordCompare(encrypted, plain string) bool {
	password := saltAes(plain)
	return nil == bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
}
