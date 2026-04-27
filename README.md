# echo-auth-playground

go-ticket-engine/
│
├── cmd/
│   └── server/
│       └── main.go                     # Wire DI: sqlc → repo → usecase → handler → router → start
│
├── db/
│   ├── migrations/
│   └── queries/
│
├── internal/
│   │
│   ├── domain/                         # ★ Thuần Go, không import bất cứ gì từ infra/delivery/framework
│   │   ├── user/
│   │   │   ├── user.go                 # Entity: struct User, constructor, business methods
│   │   │   ├── repository.go           # interface UserRepository
│   │   │   ├── service.go              # interface PasswordHasher (nếu thuộc domain logic)
│   │   │   └── errors.go
│   │   ├── refreshtoken/
│   │   │   ├── token.go
│   │   │   ├── repository.go
│   │   │   └── errors.go
│   │   └── shared/
│   │       ├── ports.go               # interface TokenGenerator, TokenHasher, IdentifierGenerator
│   │       └── errors.go              # ErrUnauthorized, ErrForbidden (shared domain errors)
│   │
│   ├── application/                   # ★ Use cases: chỉ phụ thuộc domain interface
│   │   └── auth/
│   │       ├── login/
│   │       │   ├── command.go         # struct Command (DTO input)
│   │       │   └── handler.go         # struct Handler + Execute(ctx, cmd) (*Result, error)
│   │       ├── register/
│   │       │   ├── command.go
│   │       │   └── handler.go
│   │       ├── refresh/
│   │       │   ├── command.go
│   │       │   └── handler.go
│   │       ├── logout/
│   │       │   ├── command.go
│   │       │   └── handler.go
│   │       └── me/
│   │           ├── command.go
│   │           └── handler.go
│   │
│   ├── infrastructure/                # ★ Implement interface domain
│   │   ├── persistence/
│   │   │   ├── sqlc/                  # Code gen – KHÔNG import từ domain/application
│   │   │   │   ├── db.go
│   │   │   │   ├── models.go
│   │   │   │   ├── user.sql.go
│   │   │   │   └── auth.sql.go
│   │   │   ├── user_repo.go           # implements user.UserRepository dùng sqlc
│   │   │   └── refreshtoken_repo.go   # implements refreshtoken.RefreshTokenRepository dùng sqlc
│   │   ├── crypto/
│   │   │   ├── bcrypt_hasher.go       # implements PasswordHasher
│   │   │   └── sha256_hasher.go       # implements TokenHasher
│   │   ├── token/
│   │   │   └── jwt.go                 # implements TokenGenerator
│   │   └── identifier/
│   │       └── uuid.go                # implements IdentifierGenerator
│   │
│   ├── delivery/
│   │   └── http/
│   │       ├── router.go              # khai báo toàn bộ route, inject middleware + handler
│   │       ├── response.go            # helper: OK(), Error(), Paginate()
│   │       ├── errors.go              # Global error handler + MapDomainErrorToHTTP
│   │       ├── middleware/
│   │       │   ├── auth.go            # JWT middleware
│   │       │   ├── logger.go
│   │       │   └── recover.go
│   │       └── auth/                  # ★ Tách theo module, không gom 1 file
│   │           ├── handler.go         # struct AuthHTTPHandler + inject tất cả usecase
│   │           └── routes.go          # RegisterRoutes(e *echo.Echo, h *AuthHTTPHandler)
│   │
│   ├── config/
│   │   └── config.go
│   │
│   └── database/
│       ├── postgres.go
│       └── redis.go
│
└── pkg/
    └── validator/
        └── validator.go
