package jwt

import (
	"errors"
	"fmt"
	"time"

	authmodel "github.com/arseniychernykh/test-backend-1-ChernykhITMO/internal/auth/model"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var ErrInvalidToken = errors.New("invalid token")

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Manager {
	return &Manager{
		secret: []byte(secret),
		ttl:    ttl,
	}
}

func (m *Manager) Issue(userID uuid.UUID, role authmodel.Role) (string, error) {
	const op = "internal.platform.jwt.Manager.Issue"

	now := time.Now().UTC()

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, Claims{
		UserID: userID.String(),
		Role:   string(role),
		RegisteredClaims: jwtlib.RegisteredClaims{
			IssuedAt:  jwtlib.NewNumericDate(now),
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.ttl)),
		},
	})

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", fmt.Errorf("%s: sign token: %w", op, err)
	}

	return signed, nil
}

func (m *Manager) Parse(token string) (*Claims, error) {
	const op = "internal.platform.jwt.Manager.Parse"

	parsed, err := jwtlib.ParseWithClaims(token, &Claims{}, func(t *jwtlib.Token) (any, error) {
		if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	claims, ok := parsed.Claims.(*Claims)
	if !ok || !parsed.Valid {
		return nil, fmt.Errorf("%s: %w", op, ErrInvalidToken)
	}

	return claims, nil
}
