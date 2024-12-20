package logger

const (

	// Application Messages
	ApplicationStartedMessage  = "Application started"
	ApplicationShutdownMessage = "Application shutdown"

	// Config Messages
	ConfigLoadedMessage = "config: %s"

	// Controller Messages
	SendingBatchOfPasswords = "Sending batch of passwords"
	StreamCloseErrorMessage = "Failed to close stream: %s"
	StreamClosedMessage     = "Stream closed successfully"

	// Usecase Messages
	LoadingPasswordFromFileMessage = "Loading Passwords from file"
)
