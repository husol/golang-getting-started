package handlers

import (
  "github.com/dgrijalva/jwt-go"
  "github.com/google/uuid"
  "github.com/labstack/echo"
  "golang-getting-started/models"
  "golang-getting-started/repositories"
  "net/http"
  "strconv"
)

type Result struct {
  Posts []models.Post `json:"posts"`
  TotalPage int `json:"total_page"`
}

func GetPosts() func(c echo.Context) error {
  return func(c echo.Context) error {
    repoPost := repositories.NewPostRepo()

    pageStr := c.QueryParam("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
      page = 1
    }

    limitStr := c.QueryParam("limit")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
      limit = 10
    }

    columns := []string{"id", "title", "content", "author_id"}
    var conditions []models.Condition

    posts, totalPage := repoPost.Find(columns, conditions, models.Paging{page, limit}, "")

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data: Result{
        Posts: posts,
        TotalPage: totalPage,
      },
    })
  }
}

func GetPostById() func(c echo.Context) error {
  return func(c echo.Context) error {
    id := c.Param("id")

    repoPost := repositories.NewPostRepo()
    myPost := repoPost.GetById(id)

    if myPost.ID == "" {
      return c.JSON(404, Response{
        Status:  404,
        Message: "Not Found",
      })
    }

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data: myPost,
    })
  }
}

func CreatePost() func(c echo.Context) error {
  return func(c echo.Context) error {
    paramPost := new(models.Post)
    _ = c.Bind(paramPost)

    //Validate
    if paramPost.Title == "" {
      return c.JSON(501, Response{
        Status:  501,
        Message: "Title is required.",
      })
    }

    uuid := uuid.New()
    user := c.Get("user").(*jwt.Token)
    loggedUser := user.Claims.(*models.JwtCustomClaims)

    myPost := models.Post{ID: uuid.String(), Title: paramPost.Title, Content: paramPost.Content, AuthorID: loggedUser.ID}

    repoPost := repositories.NewPostRepo()
    repoPost.Create(myPost)

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data: myPost,
    })
  }
}

func UpdatePost() func(c echo.Context) error {
  return func(c echo.Context) error {
    id := c.Param("id")

    //Check if the post existed
    repoPost := repositories.NewPostRepo()
    myPost := repoPost.GetById(id)

    if myPost.ID == "" {
      return c.JSON(404, Response{
        Status:  404,
        Message: "Not Found",
      })
    }

    user := c.Get("user").(*jwt.Token)
    loggedUser := user.Claims.(*models.JwtCustomClaims)

    if myPost.AuthorID != loggedUser.ID {
      return c.JSON(404, Response{
        Status:  403,
        Message: "Permission denied.",
      })
    }

    paramPost := new(models.Post)
    c.Bind(paramPost)

    if paramPost.Title == "" {
      return c.JSON(501, Response{
        Status:  501,
        Message: "Title is required.",
      })
    }
    myPost.Title = paramPost.Title
    myPost.Content = paramPost.Content

    //Update the post to DynamoDB
    repoPost.Update(myPost)

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data: myPost,
    })

  }
}

func RemovePost() func(c echo.Context) error {
  return func(c echo.Context) error {
    id := c.Param("id")

    if id == "" {
      return c.JSON(501, Response{
        Status:  501,
        Message: "Invalid parameters",
      })
    }

    //Check if the post existed
    repoPost := repositories.NewPostRepo()
    myPost := repoPost.GetById(id)

    if myPost.ID == "" {
      return c.JSON(501, Response{
        Status:  501,
        Message: "The post is not existed in system.",
      })
    }

    //Check if the post belong to logged user
    user := c.Get("user").(*jwt.Token)
    loggedUser := user.Claims.(*models.JwtCustomClaims)

    if myPost.AuthorID != loggedUser.ID {
      return c.JSON(404, Response{
        Status:  403,
        Message: "Permission denied.",
      })
    }

    err := repoPost.Delete(myPost)

    if err != nil {
      return c.JSON(501, Response{
        Status:  501,
        Message: err.Error(),
      })
    }

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
    })
  }
}

func GetMyPosts() func(c echo.Context) error {
  return func(c echo.Context) error {
    repoPost := repositories.NewPostRepo()

    pageStr := c.QueryParam("page")
    page, err := strconv.Atoi(pageStr)
    if err != nil {
      page = 0
    } else {
      page--
    }

    limitStr := c.QueryParam("limit")
    limit, err := strconv.Atoi(limitStr)
    if err != nil {
      limit = 10
    }

    user := c.Get("user").(*jwt.Token)
    loggedUser := user.Claims.(*models.JwtCustomClaims)

    columns := []string{"id", "title", "content", "author_id"}
    var conditions []models.Condition
    var values []interface{}
    conditions = append(conditions, models.Condition{
      Field: "author_id",
      Operator: "=",
      Values: append(values, loggedUser.ID),
    })

    posts, totalPage := repoPost.Find(columns, conditions, models.Paging{page, limit}, "")

    return c.JSON(http.StatusOK, Response{
      Status:  http.StatusOK,
      Message: "OK",
      Data: Result{
        Posts: posts,
        TotalPage: totalPage,
      },
    })
  }
}