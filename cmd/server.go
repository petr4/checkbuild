package cmd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/petr4/checkbuild/controller"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run checkbuild as server",
	Long:  "`checkbuild` works in modes: server,cli. This is server mode.",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Server Init...")
		logrus.Warnf("Server Init...:%v", viper.IsSet("server.port"))
		if !viper.IsSet("server.port") {
			logrus.Fatal("Can not open config file")
		}
		gin.SetMode(gin.ReleaseMode)
		if viper.GetBool("debug") {
			gin.SetMode(gin.DebugMode)
			gin.DisableConsoleColor()

			//logFile := fmt.Sprintf("%s/checkbuild.log", viper.GetString("log_file"))
			// f, err := os.Create(logFile)
			// if err != nil {
			// 	panic(fmt.Sprintf(
			// 		"Error while create log file [%s]: %s",
			// 		logFile,
			// 		err.Error(),
			// 	))
			// }
			// gin.DefaultWriter = io.MultiWriter(f)
		}

		r := gin.Default()
		//r.Use(middleware.Auth())
		r.GET("/", controller.Index)
		//r.GET("/_healthcheck", controller.HealthCheck)
		r.GET("/favicon.ico", func(c *gin.Context) {
			c.String(http.StatusNoContent, "")
		})
		r.Run(fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt("server.port"))))

	},
}
