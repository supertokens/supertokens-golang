package supertokens

const (
	HeaderRID                         = "rid"
	HeaderFDI                         = "fdi-version"
	ServerlessCacheBaseFilePath       = "/tmp"
	ServerlessCacheAPIVersionFilePath = ServerlessCacheBaseFilePath + "/supertokens-apiversion"
	version                           = "6.0.0"
)

var (
	cdiSupported = []string{"2.7"}
)
