/*
Copyright © 2022 Bpazy
*/
package cmd

import (
	"errors"
	"fmt"
	"github.com/Bpazy/berrors"
	"github.com/Bpazy/webhook-forwarder/model"
	"github.com/Bpazy/webhook-forwarder/util"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
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
	log.Debugf("Got template name: %s", name)
	templatesPath := getTemplatePath()
	if err := checkTemplateName(templatesPath, name); err != nil {
		return err
	}
	fileBody, err := os.ReadFile(filepath.Join(templatesPath, name))
	if err != nil {
		return err
	}
	log.Debugf("Got forwarding body: %s", string(forwardBody))
	t := util.NewTemplate(string(fileBody), forwardBody)
	runResult, err := t.RunJs()
	if err != nil {
		return err
	}
	return doRequest(runResult, err)
}

func doRequest(runResult *util.RunResult, err error) error {
	targetUrl, err := runResult.GetString("target")
	if err != nil {
		return err
	}
	payload, err := runResult.GetObject("payload")
	if err != nil {
		return err
	}
	client := resty.New()
	res, err := client.R().
		SetBody(payload).
		Post(targetUrl)
	log.Debugf("Got response: %s", res.String())
	if err != nil {
		return fmt.Errorf("forward request to %s failed: %+v", targetUrl, err)
	}
	return nil
}

func checkTemplateName(templatesPath string, name string) error {
	dirs := berrors.Unwrap(os.ReadDir(templatesPath))
	for _, dir := range dirs {
		if dir.Name() == name {
			return nil
		}
	}
	return errors.New(fmt.Sprintf("No template %s exists", name))
}

func getTemplatePath() string {
	userHomeDir := berrors.Unwrap(os.UserHomeDir())
	templatesPath := filepath.Join(userHomeDir, "/.config/webhook-forwarder/templates")
	return templatesPath
}

var port string

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&port, "port", "p", ":8080", "Serve port")
}
