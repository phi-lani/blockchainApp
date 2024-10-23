package models

import (
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/phi-lani/blockchainApp/utils"
	"gorm.io/gorm"
)

type User struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	WalletAddress string    `json:"wallet_address"`
	PrivateKey    string    `json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type MFAToken struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileUpdate struct {
	Email         string `json:"email"`
	Username      string `json:"username"`
	WalletAddress string `json:"wallet_address"`
}

func CreateUser(db *gorm.DB, user *User) error {
	user.Password = strings.TrimSpace(user.Password)
	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashedPassword

	walletAddress, privateKey, err := GenerateEthereumWallet()
	if err != nil {
		return err
	}

	user.WalletAddress = walletAddress
	user.PrivateKey = privateKey

	return db.Create(user).Error
}

func AuthenticateUser(db *gorm.DB, credentials UserCredentials) (User, error) {
	var user User
	email := strings.ToLower(credentials.Email)
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, errors.New("user not found")
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		return user, errors.New("incorrect password")
	}

	return user, nil
}

func StoreMFAToken(db *gorm.DB, userID uint, token string) error {
	mfaToken := MFAToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: time.Now().Add(10 * time.Minute),
		CreatedAt: time.Now(),
	}
	return db.Create(&mfaToken).Error
}

func UpdateUserProfile(db *gorm.DB, userID string, updates *UserProfileUpdate) error {
	return db.Model(&User{}).Where("id = ?", userID).Updates(updates).Error
}

func GenerateEthereumWallet() (string, string, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return "", "", err
	}
	privateKeyBytes := crypto.FromECDSA(privateKey)
	privateKeyHex := hex.EncodeToString(privateKeyBytes)
	publicKey := privateKey.PublicKey
	walletAddress := crypto.PubkeyToAddress(publicKey).Hex()
	return walletAddress, privateKeyHex, nil
}
