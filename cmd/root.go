package cmd

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/mritd/csi-archetype/pkg/archetype"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	Version   string
	BuildDate string
	CommitID  string

	versionTpl = `
Name: csi-archetype
Version: %s
Arch: %s
BuildDate: %s
CommitID: %s
`
)

var (
	debug    bool
	name     string
	nodeID   string
	endpoint string

	// some parameters
	parameter1 string
	parameter2 int
	parameter3 time.Duration
)

var rootCmd = &cobra.Command{
	Use:     "csi-archetype",
	Short:   "CSI archetype",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		archetype.NewCSIDriver(
			name,
			strings.TrimPrefix(Version, "v"),
			nodeID,
			endpoint,
			parameter1,
			parameter2,
			parameter3,
		).Run()
	},
}

func init() {
	// Initialize the log component after parsing command line parameters
	cobra.OnInitialize(initLog)

	// The debug option is used to control whether logrus outputs debug logs
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug log")

	// The parameters necessary for running the CSI plug-in must be set by the user
	rootCmd.PersistentFlags().StringVar(&nodeID, "nodeid", "", "csi node id")
	_ = rootCmd.MarkPersistentFlagRequired("nodeid")
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "csi endpoint")
	_ = rootCmd.MarkPersistentFlagRequired("endpoint")

	// Users generally do not need to modify the csi plugin name, so we hide the `--name` option
	rootCmd.PersistentFlags().StringVar(&name, "name", "csi-archetype", "csi name")
	_ = rootCmd.PersistentFlags().MarkHidden("name")
	_ = rootCmd.MarkPersistentFlagRequired("name")

	// The user needs to add some parameters according to actual needs to ensure
	// the successful initialization of the CSI plug-in
	rootCmd.PersistentFlags().StringVar(&parameter1, "parameter1", "", "csi parameter1")
	rootCmd.PersistentFlags().IntVar(&parameter2, "parameter2", 10, "csi parameter2")
	rootCmd.PersistentFlags().DurationVar(&parameter3, "parameter3", 10*time.Second, "csi parameter3")

	// General `--version` option output template, through which you can confirm the
	// CSI version, compilation time, git commit id and other information at runtime
	rootCmd.SetVersionTemplate(fmt.Sprintf(versionTpl, Version, runtime.GOOS+"/"+runtime.GOARCH, BuildDate, CommitID))
}

func initLog() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
