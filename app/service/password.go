package service

import "golang.org/x/crypto/bcrypt"

const bcryptCost = 10

// dummyHash dipakai saat email tidak ditemukan agar waktu respons login tetap
// mirip dengan kasus password salah (mencegah user enumeration lewat timing).
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("dummy-password"), bcryptCost)

func HashPassword(plain string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(plain), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashed), nil
}

func VerifyPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

func verifyDummyPassword(plain string) {
	_ = bcrypt.CompareHashAndPassword(dummyHash, []byte(plain))
}