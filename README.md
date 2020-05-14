# Require

* Echo framework
* DynamoDB
* JWT authentication

# Các api cần tạo

#### 1. Sign Up

**POST http://localhost:9090/users/sign_up**
  
**Request**

* name
* email
* password
  
**Response**
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYmMiOiIxMjMiLCJ4eXoiOiI0NTYifQ.4oxGUWIn-7qBmuO7sAJqh7Q6iLDAdOQ1bWexHYAw8Q4"
  }
}
```

#### 2. Sign In

**POST http://localhost:9090/users/sign_in**
  
**Request**

* email
* password
  
**Response**

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhYmMiOiIxMjMiLCJ4eXoiOiI0NTYifQ.4oxGUWIn-7qBmuO7sAJqh7Q6iLDAdOQ1bWexHYAw8Q4"
  }
}
```

Fail
```json
{
  "status": 401,
  "message": "Incorrect email or password",
  "data": null
}
```

#### 3. Create Post

**POST http://localhost:9090/posts**

**Header**
* Authorization: Token
  
**Request**

* title
* content
  
**Response**
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "id": "5392d83b-6bf1-41db-9622-e6ab0289050f",
    "title": "title123",
    "content": "content123",
    "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
  }
}
```

#### 4. Get A Post

**GET http://localhost:9090/posts/:id**

**Header**
* Authorization: Token
  
**Request**

* id
   
**Response**  

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "id": "5392d83b-6bf1-41db-9622-e6ab0289050f",
    "title": "title123",
    "content": "content123",
    "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
  }
}
```

Not found
```json
{
  "status": 404,
  "message": "Not Found",
  "data": null
}
```

#### 5. Update My Post

**PUT http://localhost:9090/posts/:id**

**Header**
* Authorization: Token
  
**Request**

* id
* title
* content
  
**Response**

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "id": "5392d83b-6bf1-41db-9622-e6ab0289050f",
    "title": "title123",
    "content": "content123",
    "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
  }
}
```

Not found
```json
{
  "status": 404,
  "message": "Not Found",
  "data": null
}
```

#### 6. Delete My Post

**DELETE http://localhost:9090/posts/:id**

**Header**
* Authorization: Token
  
**Request**

* id
  
**Response**

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": null
}
```

Not found
```json
{
  "status": 404,
  "message": "Not Found",
  "data": null
}
```

#### 7. Get All Posts (phân trang và sort by title asc)

**GET http://localhost:9090/posts?page=xx&limit=yy**

**Header**
* Authorization: Token

**Request**

* page
* limit

Nếu ko truyền, mặc định lấy page đầu tiên, limit = 10
   
**Response**  

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "posts": [
    {
        "id": "5392d83b-6bf1-41db-9622-e6ab0289050f",
        "title": "title123",
        "content": "content123",
        "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
    },
    {
        "id": "9453845e-c39b-4a51-840b-797f31c18029",
        "title": "title456",
        "content": "content456",
        "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
    }
    ]
  },
  "meta": {
    "next_page": "xx"
  }
}
```

#### 8. Get My Posts (phân trang và sort by title asc)

**GET http://localhost:9090/posts/me?page=xx&limit=yy**

**Header**
* Authorization: Token

**Request**

* page
* limit

Nếu ko truyền, mặc định lấy page đầu tiên, limit = 10
   
**Response**  

Success
```json
{
  "status": 200,
  "message": "OK",
  "data": {
    "posts": [
    {
        "id": "5392d83b-6bf1-41db-9622-e6ab0289050f",
        "title": "titleabc",
        "content": "contentabc",
        "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
    },
    {
        "id": "9453845e-c39b-4a51-840b-797f31c18029",
        "title": "titledef",
        "content": "contentdef",
        "author_id": "d96a40c5-cff4-4335-a432-8b85840390b0"
    }
    ]
  },
  "meta": {
    "next_page": "xx"
  }
}
```
