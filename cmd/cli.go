package cmd

import (
	"audio/internal/build"
	"audio/internal/config"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// ParseArgs initializes and executes the command line interface,
// returning the parsed configuration or an error.
//
// The CLI supports the following modes:
// - Interactive TUI mode (default when no command is specified)
// - One-off commands (e.g., 'list' for device listing)
// - Help and version information
func ParseArgs() (*config.Config, error) {
	buildInfo := build.GetBuildFlags()
	options := config.NewConfig()

	rootCmd := &cobra.Command{
		Use:           buildInfo.Name,
		Short:         "A simple CLI audio processing engine",
		Version:       buildInfo.Version,
		SilenceErrors: true,
		SilenceUsage:  true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd:   true,
			DisableDescriptions: true,
			DisableNoDescFlag:   true,
			HiddenDefaultCmd:    true,
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options.TUIMode = true
			return nil
		},
	}

	// Display help message
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})

	// List command
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List available audio devices and exit",
		Run: func(cmd *cobra.Command, args []string) {
			options.Command = "list"
		},
	}
	rootCmd.AddCommand(listCmd)

	// Audio Device Configuration
	rootCmd.PersistentFlags().IntVarP(&options.DeviceID, "device", "d", config.DefaultDeviceID,
		"Specify input device ID. Use 'list' command to see available devices.")
	rootCmd.PersistentFlags().IntVarP(&options.Channels, "channels", "c", config.DefaultChannels,
		"Number of channels to record (1=mono, 2=stereo)")
	rootCmd.PersistentFlags().Float64VarP(&options.SampleRate, "sample-rate", "s", config.DefaultSampleRate,
		"Sample rate, measured in Hertz (Hz)")
	rootCmd.PersistentFlags().IntVarP(&options.FramesPerBuffer, "frames-per-buffer", "b", config.DefaultFramesPerBuffer,
		"The number of frames per buffer (affects latency)")
	rootCmd.PersistentFlags().BoolVarP(&options.LowLatency, "low-latency", "l", config.DefaultLowLatency,
		"Use low latency mode for real-time processing")

	// Analysis Configuration
	rootCmd.PersistentFlags().IntVarP(&options.FFTBands, "bands", "f", config.DefaultFFTBands,
		"Number of frequency bands for FFT visualization (default: 12)")

	// Recording Configuration
	rootCmd.PersistentFlags().BoolVarP(&options.RecordInputStream, "record", "r", config.DefaultRecordInputStream,
		"Record audio from the specified input device")
	rootCmd.PersistentFlags().StringVarP(&options.OutputFile, "output", "o", config.DefaultOutputFile,
		"Output file name. Default is recording-MM-DD-YYYY-HHMMSS.wav")

	// Debug Configuration
	rootCmd.PersistentFlags().BoolVarP(&options.Verbose, "verbose", "v", config.DefaultVerbosity,
		"Show verbose output")

	// Defaults
	if options.OutputFile == "" {
		options.OutputFile = "recording-" +
			time.Now().UTC().Format("02-01-2006-150405") +
			"." + options.Format
	}

	// Execute the CLI
	rootCmd.SetArgs(os.Args[1:])
	err := rootCmd.Execute()
	if err != nil {
		return nil, err
	}

	return options, nil
}
