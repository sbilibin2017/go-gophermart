package cli

func Run(f func() error) int {
	err := f()
	if err != nil {
		return 1
	}
	return 0
}
