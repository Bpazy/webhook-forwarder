/*
Copyright © 2022 Bpazy
*/
package cmd

import (
	"fmt"
	"github.com/Bpazy/berrors"
	"github.com/Bpazy/webhook-forwarder/model"
	"github.com/Bpazy/webhook-forwarder/template"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"io"
	"net/http"
)

var port string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Serve port")
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run webhook-forwarder backend server",
	Run:   serve(),
}

func serve() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		r := gin.New()
		r.Use(gin.Recovery())
		r.Any("/forward/:name", forward)
		log.Infof("Serving on %s", port)
		berrors.Must(r.Run(port))
	}
}

func forward(c *gin.Context) {
	err := doForward(c.Param("name"), berrors.Unwrap(io.ReadAll(c.Request.Body)))
	if err != nil {
		log.Errorf("doForward error: %+v", err)
		c.JSON(http.StatusOK, model.NewFailedResult(err.Error()))
		return
	}
	c.JSON(http.StatusOK, model.NewSuccessResult(nil))
}

func doForward(name string, forwardBody []byte) error {
	log.Debugf("Requesting template name: %s, forwarding: %s", name, string(forwardBody))
	t, err := template.New(name, forwardBody)
	if err != nil {
		return err
	}
	r, err := t.RunJs()
	if err != nil {
		return err
	}
	return doRequest(r, err)
}

func doRequest(r *template.JsResult, err error) error {
	client := resty.New()
	res, err := client.R().
		SetBody(r.Payload).
		Post(r.Target)
	log.Debugf("Got forward response: %s", res.String())
	if err != nil {
		return fmt.Errorf("forward request to %s failed: %+v", r.Target, err)
	}
	return nil
}
