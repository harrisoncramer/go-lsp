package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/harrisoncramer/go-lsp/logger"
	"github.com/harrisoncramer/go-lsp/server"
	"github.com/spf13/cobra"
)

var version = "1.0.0"

func main() {
	var logPath string
	var rootDir string

	ctx := context.Background()

	var rootCmd = &cobra.Command{
		Use:   "go-lsp",
		Short: "go-lsp is a Language Server Protocol implementation",
		Run: func(cmd *cobra.Command, args []string) {
			if v, _ := cmd.Flags().GetBool("version"); v {
				fmt.Printf("go-lsp version %s\n", version)
				os.Exit(0)
			}

			logger, err := logger.NewLogger(ctx, logger.NewLoggerParams{
				LogPath: "/tmp/go-lsp.log",
				Level:   logger.DebugLevel,
			})
			if err != nil {
				log.Fatalf("failed to initialize logger: %v", err)
			}
			logger.Debug("starting server")
			s := server.NewServer(ctx, logger)
			s.Start()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&logPath, "logpath", "l", "/tmp/go-lsp.log", "Set log path")
	rootCmd.PersistentFlags().StringVarP(&rootDir, "root", "r", "", "Set root directory (default: current directory)")
	rootCmd.Flags().BoolP("version", "v", false, "Print the version number")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
