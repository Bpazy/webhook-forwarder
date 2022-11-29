/*
Copyright © 2022 Bpazy
*/
package cmd

import (
	"encoding/json"
	berrors "github.com/Bpazy/berrors"
	"github.com/gin-gonic/gin"
	"github.com/robertkrimen/otto"
	"github.com/spf13/cobra"
	"io"
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
		r.Any("/ping", ping)
		r.Any("/mirror", mirror)
		r.Any("/forward/:name", forward)
		berrors.Must(r.Run(port))
	}
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func mirror(c *gin.Context) {
	r := gin.H{}
	berrors.Must(json.Unmarshal(berrors.Unwrap(io.ReadAll(c.Request.Body)), &r))
	c.JSON(http.StatusOK, r)
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
		_ = berrors.Unwrap(vm.Run(string(berrors.Unwrap(os.ReadFile(filepath.Join(templatesPath, dir.Name()))))))

		r := gin.H{}
		requestBody := berrors.Unwrap(io.ReadAll(c.Request.Body))
		var convertValue any
		if err := json.Unmarshal(requestBody, &r); err != nil {
			convertValue = berrors.Unwrap(vm.Call("convert", nil, requestBody))
		} else {
			convertValue = berrors.Unwrap(vm.Call("convert", nil, r))
		}
		c.JSON(http.StatusOK, gin.H{
			"sourcePayload": r,
			"convertValue":  convertValue,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

var port string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Serve port")
}
