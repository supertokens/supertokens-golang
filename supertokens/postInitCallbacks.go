package supertokens

var postInitCallbacks = []func() error{}

func AddPostInitCallback(cb func() error) {
	postInitCallbacks = append(postInitCallbacks, cb)
}

func runPostInitCallbacks() error {
	for _, cb := range postInitCallbacks {
		err := cb()
		if err != nil {
			return err
		}
	}
	postInitCallbacks = []func() error{}
	return nil
}
