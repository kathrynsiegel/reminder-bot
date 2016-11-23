package helpers

// PanicIfError panics if passed an error.
func PanicIfError(err error) {
	if err != nil {
		panic(err)
	}
}
