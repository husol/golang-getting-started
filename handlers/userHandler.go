package handlers

import (
  "crypto/sha256"
  "encoding/hex"
  "github.com/dgrijalva/jwt-go"
  "github.com/google/uuid"
  "github.com/labstack/echo"
  "golang-getting-started/models"
  "golang-getting-started/repositories"
  "io"
  "net/http"
  "os"
  "time"
)

type Response struct {
  Status  int         `json:"status"`
  Message string      `json:"message"`
  Data    interface{} `json:"data"`
}
type Token struct {
  Token string `json:"token"`
}

// Create a struct to read the email and password from the request body
type Credentials struct {
  Email string `json:"email"`
  Password string `json:"password"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
  id string `json:"id"`
  name string `json:"name"`
  email string `json:"email"`
  jwt.StandardClaims
}

func SignUp() func(c echo.Context) error {
  return func(c echo.Context) error {
    userParams := new(models.User)
    err := c.Bind(userParams)

    //Check if email is already existed
    repoUser := repositories.NewUserRepo()
    myUser := repoUser.GetByEmail(userParams.Email)

    if myUser.ID != "" {
      return c.JSON(401, Response{
        Status:  401,
        Message: "The email is already used in system.",
      })
    }

    //Encrypt sha256 password
    h := sha256.New()
    io.WriteString(h, userParams.Password)
    encryptPassword := hex.EncodeToString(h.Sum(nil))

    //Create new user
    uuid := uuid.New()
    repoUser.Create(models.User{uuid.String(), userParams.Name, userParams.Email, encryptPassword})

    // Set custom claims
    claims := &models.JwtCustomClaims{
      myUser.ID,
      myUser.Name,
      myUser.Email,
      jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
      },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Generate encoded token and send it as response.
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
    if err != nil {
      return c.JSON(401, Response{
        Status:  401,
        Message: err.Error(),
      })
    }

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data:    Token{Token: tokenString},
    })
  }
}

func SignIn() func(c echo.Context) error {
  return func(c echo.Context) error {
    userParams := new(models.User)
    err := c.Bind(userParams)

    repoUser := repositories.NewUserRepo()
    myUser := repoUser.GetByEmail(userParams.Email)

    //Encrypt sha256 password
    h := sha256.New()
    io.WriteString(h, userParams.Password)
    encryptPassword := hex.EncodeToString(h.Sum(nil))

    if myUser.ID == "" || myUser.Password != encryptPassword {
      return c.JSON(401, Response{
        Status:  401,
        Message: "Incorrect email or password",
      })
    }

    // Set custom claims
    claims := &models.JwtCustomClaims{
      myUser.ID,
      myUser.Name,
      myUser.Email,
      jwt.StandardClaims{
        ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
      },
    }

    // Create token with claims
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Generate encoded token and send it as response.
    tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SIGNING_KEY")))
    if err != nil {
      return c.JSON(401, Response{
        Status:  401,
        Message: err.Error(),
      })
    }

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data:    Token{Token: tokenString},
    })
  }
}