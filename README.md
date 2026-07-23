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
DB_PORT=5433
DB_HOST=127.0.0.1

MAILTRAP_API_URL=https://send.api.mailtrap.io/api/send
MAILTRAP_API_TOKEN=your_mailtrap_api_token
MAILTRAP_FROM_EMAIL=hello@demomailtrap.co
MAILTRAP_FROM_NAME=Your Project Name
APP_BASE_URL=http://localhost:8080

JWT_SECRET=your_jwt_secret
```

Note : I'm still working on it, so all the features I mentioned above may not be there when u check. 🦥

and in the future i will be improving performance by keep working on it, for example : doing account lockout by using redis instead of database.