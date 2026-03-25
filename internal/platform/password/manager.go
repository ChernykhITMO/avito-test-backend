package password

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Manager struct {
	cost int
}

func New(cost int) *Manager {
	return &Manager{cost: cost}
}

func (m *Manager) Hash(password string) (string, error) {
	const op = "internal.platform.password.Manager.Hash"

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), m.cost)
	if err != nil {
		return "", fmt.Errorf("%s: generate hash: %w", op, err)
	}

	return string(hashed), nil
}

func (m *Manager) Compare(hash, password string) error {
	const op = "internal.platform.password.Manager.Compare"

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return fmt.Errorf("%s: compare hash: %w", op, err)
	}

	return nil
}
