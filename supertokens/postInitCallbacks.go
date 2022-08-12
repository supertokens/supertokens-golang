package supertokens

var postInitCallbacks = []func(){}

func AddPostInitCallback(cb func()) {
	postInitCallbacks = append(postInitCallbacks, cb)
}

func runPostInitCallbacks() {
	for _, cb := range postInitCallbacks {
		cb()
	}
	postInitCallbacks = []func(){}
}
