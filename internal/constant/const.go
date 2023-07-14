package constant

const (
	AuthHeader    = "Authorization"
	AuthPrefix    = "Bearer"
	TokenLength   = 16
	CommitLength  = 64
	UserIDKey     = "user_id"
	DefaultBranch = "main"
	FileSavaDir   = "blobs"
)

const (
	MinUserNameLength = 1
	MaxUserNameLength = 200
	UserNamePattern   = "^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$"

	MinPasswordLength = 6
	MaxPasswordLength = 50
	PasswordPattern   = "[a-zA-Z0-9~!@&%#_]"

	MinRepositoryNameLength = 1
	MaxRepositoryNameLength = 200
	RepositoryNamePattern   = "^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$"

	MinDraftLength = 1
	MaxDraftLength = 20
	DraftPattern   = "^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$"

	MinTagLength = 1
	MaxTagLength = 20
	TagPattern   = "^[a-zA-Z][a-zA-Z0-9_-]*[a-zA-Z0-9]$"

	MinPageSize = 1
	MaxPageSize = 50
)
