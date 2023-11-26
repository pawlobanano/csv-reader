package customerimporter

// Customer struct represents a customer record in the CSV file.
type Customer struct {
	FirstName string
	LastName  string
	Email     string
	Gender    string
	IPAddress string
}

// Logger is an interface representing the required logging methods.
type Logger interface {
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

// MockLogger is a mock implementation of the slog.Logger interface for testing.
type MockLogger struct {
	InfoCalled  bool
	WarnCalled  bool
	ErrorCalled bool
}

// Task is a struct of CSV file records (size of this slice depends of the reader's buffer).
type Task struct {
	record []string
}

// DomainCounter is a struct made for convenience for the results channel.
type DomainCounter struct {
	domain  string
	counter int
}

// Config is a struct which encapsulates .env file variables.
type Config struct {
	Concurrency              int
	InputCSVFilePathDefault  string
	InputCSVFilePath0Lines   string
	InputCSVFilePath10Lines  string
	InputCSVFilePath3kLines  string
	InputCSVFilePath10mLines string
	ReadBufferSizeInBytes    int
}
