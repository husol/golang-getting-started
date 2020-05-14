package main

import (
	"github.com/labstack/echo"
	"net/http"
)

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Posts struct {
	Posts []Post `json:"posts"`
}

type Post struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID string `json:"author_id"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

func main() {
	e := echo.New()
	e.POST("/users/sign_up", signUp())
	e.POST("/users/sign_in", signIn())
	e.POST("/posts", postFunc())
	e.GET("/posts/:id", postWithID())
	e.PUT("/posts/:id", putWithID())

	i := 0
	e.DELETE("/posts/:id", deleteWithID(&i))

	e.GET("/posts", getPost())

	e.GET("/posts/me", getPostCuaTui())

	e.Logger.Fatal(e.Start(":9090"))
}

func getPostCuaTui() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{
			Status:  http.StatusOK,
			Message: "OK",
			Data: Posts{
				Posts: []Post{
					{
						ID:       "my post id 2",
						Title:    "title2",
						Content:  "content2",
						AuthorID: "authorId2",
					},
					{
						ID:       "my post id 3",
						Title:    "title3",
						Content:  "content3",
						AuthorID: "authorId3",
					},
				},
			},
		})
	}
}

func getPost() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{
			Status:  http.StatusOK,
			Message: "OK",
			Data: Posts{
				Posts: []Post{
					{
						ID:       "my post id 1",
						Title:    "title1",
						Content:  "content1",
						AuthorID: "authorId1",
					},
					{
						ID:       "my post id 2",
						Title:    "title2",
						Content:  "content2",
						AuthorID: "authorId2",
					},
					{
						ID:       "my post id 3",
						Title:    "title3",
						Content:  "content3",
						AuthorID: "authorId3",
					},
				},
			},
		})
	}
}

func deleteWithID(i *int) func(c echo.Context) error {
	return func(c echo.Context) error {
		if *i == 0 {
			*i++
			return c.JSON(http.StatusOK, Response{
				Status:  http.StatusOK,
				Message: "OK",
			})
		}
		return c.JSON(404, Response{
			Status:  404,
			Message: "Not Found",
		})
	}
}

func putWithID() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		post := new(Post)
		_ = c.Bind(post)
		if id == "my post id" {
			return c.JSON(http.StatusOK, Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: Post{
					ID:       "my post id",
					Title:    post.Title,
					Content:  post.Content,
					AuthorID: "authorId1",
				},
			})
		}
		return c.JSON(404, Response{
			Status:  404,
			Message: "Not Found",
		})
	}
}

func postWithID() func(c echo.Context) error {
	return func(c echo.Context) error {
		id := c.Param("id")
		if id == "my post id" {
			return c.JSON(http.StatusOK, Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data: Post{
					ID:       "my post id",
					Title:    "title1",
					Content:  "content1",
					AuthorID: "authorId1",
				},
			})
		}
		return c.JSON(404, Response{
			Status:  404,
			Message: "Not Found",
		})
	}
}

func postFunc() func(c echo.Context) error {
	return func(c echo.Context) error {
		post := new(Post)
		_ = c.Bind(post)
		return c.JSON(http.StatusOK, Response{
			Status:  http.StatusOK,
			Message: "OK",
			Data: Post{
				ID:       "my post id",
				Title:    post.Title,
				Content:  post.Content,
				AuthorID: "authorId1",
			},
		})
	}
}

func signIn() func(c echo.Context) error {
	return func(c echo.Context) error {
		user := new(User)
		_ = c.Bind(user)
		if user.Password == "password1" {
			return c.JSON(http.StatusOK, Response{
				Status:  http.StatusOK,
				Message: "OK",
				Data:    Token{Token: "my token"},
			})
		}
		return c.JSON(401, Response{
			Status:  401,
			Message: "Incorrect email or password",
		})
	}
}

func signUp() func(c echo.Context) error {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, Response{
			Status:  http.StatusOK,
			Message: "OK",
			Data:    Token{Token: "token1"},
		})
	}
}
