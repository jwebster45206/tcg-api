package auth

type Claims struct {
	Sub    string   `json:"sub"`     // Subject (user ID or principal identifier)
	Iss    string   `json:"iss"`     // Issuer (token issuer identifier)
	Aud    string   `json:"aud"`     // Audience (intended recipient or service)
	Iat    int64    `json:"iat"`     // Issued At (Unix seconds when token was created)
	Exp    int64    `json:"exp"`     // Expiration Time (Unix seconds when token expires)
	Scopes []string `json:"scopes"`  // Scopes (fine-grained permissions)
	Roles  []string `json:"roles"`   // Roles (coarse-grained access roles)
}
