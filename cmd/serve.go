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
	tem "html/template"
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

var homePageHtml = tem.Must(tem.New("homePageHtml").Parse(`
<html>
<head>
  <title>webhook-forwarder</title>
</head>
<body>
  <span>Here is <a href="https://github.com/Bpazy/webhook-forwarder">webhook-forwarder</a></span>
</body>
</html>
`))

func serve() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		r := gin.New()
		r.Use(gin.Recovery())

		r.SetHTMLTemplate(homePageHtml)
		r.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "homePageHtml", nil)
		})
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
	content, err := template.GetTemplateContent(name)
	if err != nil {
		return err
	}
	t, err := template.New(content, forwardBody)
	if err != nil {
		return err
	}
	r, err := t.RunJs()
	if err != nil {
		return err
	}
	return doRequest(r)
}

func doRequest(r *template.JsResult) error {
	client := resty.New()
	for _, target := range r.Targets {
		res, err := client.R().
			SetBody(r.Payload).
			Post(target)
		log.Debugf("Got forward response: %s", res.String())
		if err != nil {
			return fmt.Errorf("forward request to %s failed: %+v", target, err)
		}
	}
	return nil
}
