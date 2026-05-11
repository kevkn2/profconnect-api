package port

type PasswordHasher interface {
	Hash(plaintext string) (string, error)
	Verify(hashed, plaintext string) error
}