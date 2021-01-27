package controller

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/petr4/checkbuild/pkg/cmp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Data struct {
	Results []cmp.Result
	Ok      bool
}

// Index controller
func Index(c *gin.Context) {
	var results []cmp.Result
	var ok bool
	//ss, _ := cmd.Flags().GetStringSlice("urls")
	ss1 := viper.GetStringSlice("urls")
	var ss []string
	if len(ss1) >= 1 {
		ss = append(ss1[:0:0], ss1...)
	}
	// debug, _ := cmd.Flags().GetBool("debug")
	// logfile, _ := cmd.Flags().GetBool("debug")
	if viper.GetBool("debug") {
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.Infof("Urls: %v", ss)
	cmpcl, err := cmp.Init()
	if err != nil {
		logrus.Fatalf("Can not init http, err %v", err)
	}
	results, ok, err = cmpcl.Run(ss)
	if err != nil {
		logrus.Info("Error: %v\n", err)

	}
	if ok {
		for _, r := range results {
			logrus.Infof("%v:%v: True\n", r.Url, r.Build)
		}
		logrus.Info("Test: passed")
	} else {
		for _, r := range results {
			logrus.Infof("%v:%v: False\n", r.Url, r.Build)
		}

		logrus.Info("Test: failed")

	}

	data := Data{
		Results: results,
		Ok:      ok}
	// Render a template with our page data
	tmpl, err := htmlTemplate(data)

	// If we got an error, write it out and exit
	if err != nil {

		c.Writer.WriteHeader(http.StatusBadRequest)
		c.Writer.Write([]byte(err.Error()))
		return
	}

	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write([]byte(tmpl))
	return
}

func htmlTemplate(pd Data) (string, error) {
	// Define a basic HTML template
	html := `<!doctype html><html> <head> <meta charset="utf-8"> <meta name="viewport" content="width=device-width, initial-scale=1"> <title>checkbuild</title> <link href="https://fonts.googleapis.com/css?family=Nunito:200,600" rel="stylesheet" type="text/css"> <style> html, body { background-color: #fff; color: #636b6f; font-family: 'Nunito', sans-serif; font-weight: 200; height: 100vh; margin: 0; } .full-height { height: 100vh; } .flex-center { align-items: center; display: flex; justify-content: center; } .position-ref { position: relative; } .top-right { position: absolute; right: 10px; top: 18px; } .content { text-align: center; } .title { font-size: 30px; } .links > a { color: #636b6f; padding: 0 25px; font-size: 12px; font-weight: 600; letter-spacing: .1rem; text-decoration: none; text-transform: uppercase; } .m-b-md { margin-bottom: 30px; } </style> </head>
	<body>
	  <div class="flex-center position-ref full-height">
		<div class="content">
		  <div class="title"> Checkbuild </div>
			<p>Test results: <span style="font-size: 150%;font-weight:bold; background-color:#DDFF33;">&nbsp;{{.Ok}}&nbsp;</span></p>
			<div>
			<ul style="display: inline-block; text-align: left;">
			{{range $result := .Results}}
			  <li><span style="color:black">{{$result.Url}}: {{$result.Build}} : {{$.Ok}}</span></li>
			{{end}}
             </ul>
			</div>
		</div> 
	   </div> 
	</body></html>`

	// Parse the template
	tmpl, err := template.New("index").Parse(html)
	if err != nil {
		return "", err
	}

	// We need somewhere to write the executed template to
	var out bytes.Buffer

	// Render the template with the data we passed in
	if err := tmpl.Execute(&out, pd); err != nil {
		// If we couldn't render, return a error
		return "", err
	}

	// Return the template
	return out.String(), nil
}
