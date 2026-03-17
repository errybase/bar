package bar

func MayPanic(err error) {
	if err != nil {
		panic(err)
	}
}
