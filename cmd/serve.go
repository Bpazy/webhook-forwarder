/*
Copyright © 2022 Bpazy
*/
package cmd

import (
	"github.com/Bpazy/berrors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"net/http"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run webhook-forwarder backend server",
	Run:   serve(),
}

func serve() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		r := gin.Default()
		r.GET("/ping", ping)
		r.GET("/forward/:name", forward)
		berrors.Must(r.Run(port))
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func forward(c *gin.Context) {
	_ = c.Param("name")
}

var port string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Serve port")
}
