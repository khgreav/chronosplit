// Package version provides exports for project version and build date.
package version

import "fmt"

var (
	Version   = "dev"
	BuildDate = "unknown"
)

func String() string {
	return fmt.Sprintf("Version %s (built %s)", Version, BuildDate)
}
