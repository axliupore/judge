package status

// Define various possible statuses representing program execution results in an online judging system.
const (
	Accepted            = "Accepted"              // Program executed successfully and provided correct output
	MemoryLimitExceeded = "Memory Limit Exceeded" // Program exceeded allowed memory limit
	TimeLimitExceeded   = "Time Limit Exceeded"   // Program exceeded allowed time limit
	OutputLimitExceeded = "Output Limit Exceeded" // Program generated output exceeding limit
	FileError           = "File Error"            // Error occurred while handling files
	NonzeroExitStatus   = "Nonzero Exit Status"   // Program ended with a non-zero exit code
	Signalled           = "Signalled"             // Program was interrupted by a signal
	InternalError       = "Internal Error"        // Internal system error
	ParamsError         = "Params Error"          // Invalid parameters passed
	ServerError         = "Server Error"          // Server error
)
