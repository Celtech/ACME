# Authorization

To authenticate with this API you need to get a valid JWT token from your credentials.
This JWT token should then be passed in the `Authorization:` header as a bearer token.

```text
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJuYW1lIjoiMTIzNDU2Nzg5IiwidXNlciI6dHJ1ZSwiZXhwIjoxNjU0NzE1NjMzLCJpYXQiOjE2NTQ1NDI4MzMsImlzcyI6IlJ5a2VMYWJzIn0.j4TH9NhImar-rj4VeNdqNMCILW3qVEg-XXltFYyZgs8
```

JWT tokens have a life span of 1 hour and must be re-obtained on expiration.

See the [token endpoint documentation](#tag/Token/paths/~1token/post) for
more information on obtaining a JWT token.

<!-- ReDoc-Inject: <security-definitions> -->
