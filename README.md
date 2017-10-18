## Summary

This is a demo project that shows the use of JSON Web Tokens (https://jwt.io/).

The JWT tokens are created with a `auth_token` in the claims.   A logged in user
sends a header with the JWT token.  The signature of the token is validated
and then the `auth_token` is extracted and validated in Postgres.

## API Specs

### `POST /signup`
Endpoint to register a new user.

```json
{
  "email": "test@example.com",
  "password": "aPassword",
  "firstName": "FirstName",
  "lastName": "LastName"
}
```

`email` will be used as a unique key in the database.

All responses are in JSON.  This one will include the "token" field with the JWT token in it.

```json
{
  "token": "some_jwt_token" 
}
```

### `POST /login`
Login a user endpoint.  The body will have the following:

```json
{
  "email": "test@example.com",
  "password": "aPassword"
}
```

Like register the response is the JWT token.

```json
{
  "token": "some_jwt_token"
}
```

### `GET /users`
Retrieve a JSON object with all of the users in it. 
The user must be logged in with a valid `X-Authentication-Token` header.

```json
{
  "users": [
    {
      "email": "test@example.com",
      "firstName": "FirstName",
      "lastName": "LastName"
    }
  ]
}
```

### `PUT /users`

This will update the user table.  Only two fields are updated.  The user must be logged into to 
perform the update.  The update will be for the user's own data.  The user will be identified
by the `auth_token`.

```json
{
  "firstName": "NewFirstName",
  "lastName": "NewLastName"
}
```

The response can body will be

```json
{
  "status": "success",
}
```

on a successful update.

## Instalation instructions.

Look at INSTRUCTIONS.md.

