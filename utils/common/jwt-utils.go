package common

import (
	"errors"

	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/config"
	"github.com/Ak8388/applikasi-antrian-dokter-gigi-payment-service/model"
	"github.com/golang-jwt/jwt/v5"
)

type JwtToken interface {
	VerfifyToken(model model.TokenAkses) (jwt.MapClaims, error)
}

type jwtToken struct {
	cfg *config.Config
}

func (j *jwtToken) VerfifyToken(model model.TokenAkses) (jwt.MapClaims, error) {
	token, err := jwt.Parse(model.TokenString, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.GetSigningMethod("HS256") {
			return nil, errors.New("token methode not match")
		}

		return j.cfg.SecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)

	if !token.Valid || !ok {
		return nil, errors.New("token not valid")
	}

	return mapClaims, nil
}

func NewJwtUtils(cfg *config.Config) JwtToken {
	return &jwtToken{cfg}
}
