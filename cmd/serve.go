/*
Copyright © 2022 Bpazy
*/
package cmd

import (
	berrors "github.com/Bpazy/berrors"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"path/filepath"
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
	name := c.Param("name")
	userHomeDir := berrors.Unwrap(os.UserHomeDir())
	templatesPath := filepath.Join(userHomeDir, "/.config/webhook-forwarder/templates")
	dirs := berrors.Unwrap(os.ReadDir(templatesPath))
	for _, dir := range dirs {
		if dir.Name() != name {
			continue
		}

		vm := otto.New()
		b := berrors.Unwrap(os.ReadFile(filepath.Join(templatesPath, dir.Name())))
		v := berrors.Unwrap(vm.Run(string(b)))
		c.JSON(http.StatusOK, v)
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

var port string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Serve port")
}
