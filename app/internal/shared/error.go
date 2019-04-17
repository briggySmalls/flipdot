package shared

// In-queue error handler
func ErrorHandler(err error) {
	if err != nil {
		panic(err)
	}
}
