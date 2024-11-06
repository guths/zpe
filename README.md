# ZPE Code Challenge
The language choice was golang due to the ease and simplicity of development.

## Extras 

- GitHub Actions to run the tests

## How to run

#### To up the api

```
docker compose up app
```

inside the container, to populate users run

```
go run seeds/main.go --table=user_role
```
__________________

#### To run the tests 

```
docker compose run --rm test
```


## Stack 
| Tech      | Version |
| ---------      | -------|
| Golang         | 1.23.2   |
| Docker         | 27.0.3 |
| Docker compose | v2.28.1 |

## How I solved the challenge

- Firstly, I thought of three tables, users, roles and user_roles, as the statement said that a user could have more than one role, I went that way.
    - I'm using ORM Gorm, so using auto migrates helps me a lot
- The route and resource access logic was designed as follows
    - First, I created a seed to have three users each with a permission level
        - **The seed can be executed like that: go run seeds/main.go --table=user_role**
    - The user can login, the auth route generate a stateless jwt token that gonna be used in other routes
    - When the user is authenticated, is possible to:
        - Create another user
        - Update an user
        - Delete an user
    - All the routes beyond the Login route are protected routes, the routes that modify data need modifier role, and the routes to just see information needs watcher role.
    - The user cannot create another user who has greater permissions than theirs
    - The user cannot delete a user with higher permission than theirs
    - The user cannot update another user permissions to be higher than theirs

### Tests

Due to time, only integration tests were developed, as they cover a larger part of the code

## Routes 

### Authentication

POST api/v1/auth/login
- Description: Authenticates a user and returns a JWT token.
- Headers:
- Content-Type: application/json
Request Body:
```
{
  "email": "user@example.com",
  "password": "password"
}
```

Response:
- 200 OK on success with JWT token in the response body.
- 401 Unauthorized on failure.

### User management

**POST /api/v1/user**
- Description:  Creates a new user. Only accessible by users with level 2 (modifier) role.
- Headers:
    - ```Authorization: Bearer <token>```
- Content-Type: application/json
Request Body:
```
{
    "username": "jonathan",
    "email": "guths@admin.com",
    "password": "guths",
    "roles": ["modifier", "watcher"]
}
```

Response:
- 200 OK on success with JWT token in the response body.
- 401 Unauthorized on failure.
- 409 User already exists.


**GET /api/v1/user/:email**
- Description:  Retrieves user details by email. Accessible by all authenticated users.
- Headers:
    - ```Authorization: Bearer <token>```
- Content-Type: application/json

Response:
- 200 OK on success with JWT token in the response body.
- 401 Unauthorized on failure.
- 404 Not Found if the user does not exist.


**DELETE /api/v1/user/:email**
- Description: Deletes a user by email. Only accessible by users with level 2 (modifier) role.
- Headers:
    - ```Authorization: Bearer <token>```
- Content-Type: application/json

Response:
- 200 OK on success.
- 401 Forbidden if the user does not have the required role.
- 404 Not Found if the user does not exist.

**PUT /api/v1/user/:email**
-  Updates user details by email. Only accessible by users with level 2 (modifier) role.
- Headers:
    - ```Authorization: Bearer <token>```
- Content-Type: application/json

Response:
- 200 OK on success.
- 401 Forbidden if the user does not have the required role.
- 404 Not Found if the user does not exist.
Request Body:
```
{
    "username": "jonathan",
    "roles": ["modifier", "watcher"]
}
```




