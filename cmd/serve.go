package cmd

import (
	"context"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/spf13/cobra"

	"github.com/riverphillips/envoy-control-plane/internal/logger"
	"github.com/riverphillips/envoy-control-plane/internal/processor"
	"github.com/riverphillips/envoy-control-plane/internal/server"
	"github.com/riverphillips/envoy-control-plane/internal/watcher"
)

var (
	port                   *uint
	watchDirectoryFileName *string
	config                 string
	nodeId                 string
	serveCmd               = &cobra.Command{
		Use:   "serve",
		Short: "Serves the xDS server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {

			cache := cache.NewSnapshotCache(false, cache.IDHash{}, logger.New(log))

			proc := processor.New(cache, nodeId, log)

			proc.ProcessFile(watcher.NotifyMessage{
				Operation: watcher.Create,
				FilePath:  *watchDirectoryFileName,
			})

			notifyCh := make(chan watcher.NotifyMessage)

			go func() {
				watcher.Watch(*watchDirectoryFileName, notifyCh)
			}()

			go func() {
				ctx := context.Background()
				srv := serverv3.NewServer(ctx, cache, nil)
				server.Serve(ctx, srv, *port)
			}()

			for {
				select {
				case msg := <-notifyCh:
					proc.ProcessFile(msg)
				}
			}
		},
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().UintVarP(port, "port", "p", 8080, "Port to listen on")
	serveCmd.Flags().StringVarP(
		watchDirectoryFileName,
		"watchDirectoryFileName",
		"wd",
		"config/config.yaml",
		"full path to directory to watch",
	)

}
