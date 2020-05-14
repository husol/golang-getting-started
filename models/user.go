package models

import "github.com/dgrijalva/jwt-go"

// JwtCustomClaims are custom claims extending default ones.
type JwtCustomClaims struct {
  ID  string `json:"id"`
  Name string `json:"name"`
  Email string   `json:"email"`
  jwt.StandardClaims
}

type User struct {
  ID       string `json:"id"`
  Name     string `json:"name"`
  Email    string `json:"email"`
  Password string `json:"password"`
}
