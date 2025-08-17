package auth

import (
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenIssuer interface {
	Issue(sub string, scopes []string, ttl time.Duration, roles []string) (string, Claims, error)
}

type TokenVerifier interface {
	Verify(token string) (Claims, error)
}

type HS256IssuerVerifier struct {
	Issuer     string
	Audience   string
	CurrentKID string
	Keys       map[string]string
}

func (h HS256IssuerVerifier) Issue(sub string, scopes []string, ttl time.Duration, roles []string) (string, Claims, error) {
	secret, ok := h.Keys[h.CurrentKID]
	if !ok {
		return "", Claims{}, errors.New("current key not found")
	}
	now := time.Now().Unix()
	exp := time.Now().Add(ttl).Unix()
	claims := Claims{
		Sub:    sub,
		Iss:    h.Issuer,
		Aud:    h.Audience,
		Iat:    now,
		Exp:    exp,
		Scopes: scopes,
		Roles:  roles,
	}
	jwtClaims := jwt.MapClaims{
		"sub":     claims.Sub,
		"iss":     claims.Iss,
		"aud":     claims.Aud,
		"iat":     claims.Iat,
		"exp":     claims.Exp,
		"scopes":  claims.Scopes,
		"roles":   claims.Roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	token.Header["kid"] = h.CurrentKID
	var keyBytes []byte
	if dec, err := base64.StdEncoding.DecodeString(secret); err == nil {
		keyBytes = dec
	} else {
		keyBytes = []byte(secret)
	}
	signed, err := token.SignedString(keyBytes)
	if err != nil {
		return "", Claims{}, err
	}
	return signed, claims, nil
}

func (h HS256IssuerVerifier) Verify(tokenStr string) (Claims, error) {
	parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		kid, _ := t.Header["kid"].(string)
		if kid == "" {
			kid = h.CurrentKID
		}
		secret, ok := h.Keys[kid]
		if !ok {
			return nil, errors.New("key not found")
		}
		if dec, err := base64.StdEncoding.DecodeString(secret); err == nil {
			return dec, nil
		}
		return []byte(secret), nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !parsed.Valid {
		return Claims{}, errors.New("invalid token")
	}
	mc, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return Claims{}, errors.New("invalid claims")
	}
	iss, _ := mc["iss"].(string)
	aud, _ := mc["aud"].(string)
	if h.Issuer != "" && iss != h.Issuer {
		return Claims{}, errors.New("issuer mismatch")
	}
	if h.Audience != "" && aud != h.Audience {
		return Claims{}, errors.New("audience mismatch")
	}
	scopes := []string{}
	if s, ok := mc["scopes"].([]interface{}); ok {
		for _, v := range s {
			if str, ok := v.(string); ok {
				scopes = append(scopes, str)
			}
		}
	}
	roles := []string{}
	if r, ok := mc["roles"].([]interface{}); ok {
		for _, v := range r {
			if str, ok := v.(string); ok {
				roles = append(roles, str)
			}
		}
	}
	c := Claims{
		Sub:    mc["sub"].(string),
		Iss:    iss,
		Aud:    aud,
		Iat:    int64(mc["iat"].(float64)),
		Exp:    int64(mc["exp"].(float64)),
		Scopes: scopes,
		Roles:  roles,
	}
	return c, nil
}
