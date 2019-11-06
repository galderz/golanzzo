package util

// ExpectNoError checks if an error exists, and if so halts execution
func ExpectNoError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
