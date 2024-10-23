package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

// Generate a 6-digit MFA token
func GenerateMFAToken() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000)) // Returns a 6-digit random number
}

func SendMFAToken(email string, token string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@example.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Your MFA Token")
	m.SetBody("text/plain", "Your MFA token is: "+token)

	d := gomail.NewDialer("smtp.gmail.com", 587, "philaningumede@gmail.com", "bpxlrlyjjdblyame")
	err := d.DialAndSend(m)
	return err
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func main() {
	// mfaToken := GenerateMFAToken()
	// err := SendMFAToken("217021035@stu.ukzn.ac.za", mfaToken)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	hshdpass, err := HashPassword("#Pn19970104!")
	if err != nil {
		log.Println(err)
	}
	fmt.Println(hshdpass)
	fmt.Println(!CheckPasswordHash("#Pn19970104!", string(hshdpass)))
}
