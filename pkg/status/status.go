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
)

// statusToCode maps status strings to integer codes
var statusToCode = map[string]int{
	Accepted:            0,
	MemoryLimitExceeded: 1,
	TimeLimitExceeded:   2,
	OutputLimitExceeded: 3,
	FileError:           4,
	NonzeroExitStatus:   5,
	Signalled:           6,
	InternalError:       7,
	ParamsError:         8,
}

// codeToStatus maps integer codes to status strings
var codeToStatus = map[int]string{
	0: Accepted,
	1: MemoryLimitExceeded,
	2: TimeLimitExceeded,
	3: OutputLimitExceeded,
	4: FileError,
	5: NonzeroExitStatus,
	6: Signalled,
	7: InternalError,
	8: ParamsError,
}

// Code function returns the integer code corresponding to the given status string.
// Different status codes can be used in the system to represent different types of errors or execution results.
func Code(status string) int {
	if code, exists := statusToCode[status]; exists {
		return code
	}
	return 7 // Unknown status defaults to 7, same as InternalError
}

// Status function returns the status string corresponding to the given integer code.
func Status(code int) string {
	if status, exists := codeToStatus[code]; exists {
		return status
	}
	return InternalError // Unknown code defaults to InternalError
}
