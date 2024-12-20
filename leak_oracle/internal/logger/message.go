package logger

const (

	// Application Messages
	ApplicationStartedMessage  = "Application started"
	ApplicationShutdownMessage = "Application shutdown"

	// Config Messages
	ConfigLoadedMessage = "config: %s"

	// Database Messages
	DatabaseConnectedMessage     = "Connected to database"
	DatabaseClosedMessage        = "Closed database connection"
	StoreBatchOfPasswordsMessage = "Stored batch of passwords"
	HashAlreadyExistsMessage     = "Password Hash already exists"
	DeadlockDetectedMessage      = "Deadlock detected, Retrying"
	TooManyConnectionsMessage    = "Too many connections, Retrying"

	// Controller Messages
	RunGRPCServerMessage       = "Starting gRPC server"
	ShutdownGRPCServerMessage  = "Shutting down gRPC server"
	RunHTTPServerMessage       = "Starting HTTP server"
	ShutdownHTTPServerMessage  = "Shutting down HTTP server"
	OpeningStream              = "Opening batch upload stream"
	ReceivingBatchOfPasswords  = "Receiving batch of passwords"
	UploadCompleteWithErrors   = "Upload complete with errors"
	UploadCompleteSuccessfully = "Upload complete"

	// UseCase Messages
	ProcessBatchPasswordsMessage    = "Processing batch of passwords"
	HashBatchOfPasswordsMessage     = "Hashing batch of passwords"
	CheckingPasswordPresenceMessage = "Checking password presence"
)
