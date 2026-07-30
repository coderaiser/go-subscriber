package meta

// These variables are set by the root package via flags.SetMeta()
// before Parse() is called, avoiding import cycles.
var (
	versionLine func() string
	help        func() string
)

// Set is called once at startup from cmd/subscriber/server.go.
func Set(vl func() string, h func() string) {
	versionLine = vl
	help = h
}

func VersionLine() string {
	if versionLine == nil {
		return "v0.0.0"
	}
	return versionLine()
}

func Help() string {
	if help == nil {
		return "usage: subscriber [options]\n"
	}
	return help()
}
