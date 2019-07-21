package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// these variables should be set at compile time as:
// go build -ldflags "-X 'github.com/mariusgiger/iot-abe/cmd.BuildTime=${$(date -u '+%Y-%m-%dT%H:%M:%SZ')}' -X 'github.com/mariusgiger/iot-abe/cmd.GitHash=${$(git log -1 --format='%H')}' -X 'github.com/mariusgiger/iot-abe/cmd.Version=${$(git describe --tags 2>/dev/null)}'" -o output/iot-abe
var (
	// Version is the service version.
	Version = "0.0.0"

	// GitHash is the hash of git commit the service is built from.
	GitHash = "unknown"

	// BuildTime build time in RFC3339 format
	BuildTime = "unknown"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of iot-abe",
	Long:  `All software has versions. This is iot-abe's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %v\n", Version)
		fmt.Printf("GitHash: %v\n", GitHash)
		fmt.Printf("BuildTime: %v\n", BuildTime)
	},
}
