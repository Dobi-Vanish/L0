package errormsg

import "errors"

// ErrorResponse represents standard error response structure.
type ErrorResponse struct {
	Error   bool   `json:"error" example:"true"`
	Message string `json:"message" example:"error description"`
}

var (
	ErrInvalidGUID                   = errors.New("invalid GUID format")
	ErrPasswordLength                = errors.New("password must be at least 8 characters long")
	ErrUserNotFound                  = errors.New("user does not exist")
	ErrInvalidID                     = errors.New("provided ID is invalid")
	ErrInvalidToken                  = errors.New("invalid access token from errormsg")
	ErrEmptyGUID                     = errors.New("empty GUID parameter")
	ErrUnexpectedSigningMethod       = errors.New("unexpected signing method")
	ErrTokenValidation               = errors.New("token validation failed")
	ErrApplyMigrations               = errors.New("error during applying migrations")
	ErrConnectDB                     = errors.New("error during connecting to DB")
	ErrSetDialect                    = errors.New("error during setting dialect to postgres")
	ErrJSONDecode                    = errors.New("JSON decode has failed")
	ErrJSONMustContain               = errors.New("must contain at least one JSON value")
	ErrDSNRequired                   = errors.New("DSN is required")
	ErrServerPortRequired            = errors.New("server port is required")
	ErrPostgresConnectAttemptsFailed = errors.New("failed connect to Postgres after 10 attempts")
	ErrTokenExpired                  = errors.New("your auth has expired, please, authenticate again")
	ErrInvalidIP                     = errors.New("invalid IP while working with tokens")
	ErrCompareHash                   = errors.New("error during comparing hash and sotre token")
	ErrPairIDNotValid                = errors.New("pair id is not valid")
	ErrGenerateRefreshToken          = errors.New("error during generating refresh token")
	ErrGenerateAccessToken           = errors.New("error during generating access token")
	ErrDeleteRefreshToken            = errors.New("error occurred during deleting refresh token")
	ErrExtractGUIDFromToken          = errors.New("error occurred during extracting guid from token")
)
