A production ready authentication backend structure uses postgresql & gorm written in go.

# Includes : 

- Login & Register flow.
- Dependency Injection with interfaces.
- Postgresql database with gorm.
- Password hashing with bcrypt.
- Running app & db with docker compose.
- JWT.
- Refresh token rotation.
- Custom error messages.
- Email verification.
- Two-factor authentication with email.
- Rate limiting.
- Request body validation.
- Structured logging.
- Health check endpoint.
- Unit tests.
- CI/CD pipeline.

# .env file example
```.env
DB_USER=GoProdReadyAuthUser
DB_PASSWORD=GoProdReadyAuth13579**
DB_NAME=GoProdReadyAuth
DB_PORT=5432
```

Note : I'm still working on it, so all the features I mentioned above may not be there when u check. 🦥