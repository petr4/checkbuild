package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/petr4/checkbuild/pkg/cmp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	CfgFile     string
	logFile     string
	userLicense string
	urls        []string
	debug       bool

	rootCmd = &cobra.Command{
		Use:   "checkbuild",
		Short: "'checkbuild' apps check build version and return success if check passed",
		Long: `'checkbuild' apps check status, compare build version in URLs, Both contain a build number(git hash).
Make sure they are the same. If they are the same, test passes. If they are not, test fails.
This will require doing some light parsing.
`,
		Run: func(cmd *cobra.Command, args []string) {
			//flag.StringVar()
			var results []cmp.Result
			var ok bool
			ss, _ := cmd.Flags().GetStringSlice("urls")
			ss1 := viper.GetStringSlice("urls")
			if len(ss1) >= 1 {
				ss = append(ss1[:0:0], ss1...)
			}
			debug, _ := cmd.Flags().GetBool("debug")
			logfile, _ := cmd.Flags().GetBool("debug")
			if debug || viper.GetBool("debug") {
				logrus.SetLevel(logrus.InfoLevel)
			}
			logrus.Infof("Urls: %v", ss)
			logrus.Infof("Debug: %v", debug)
			logrus.Infof("Logfile: %v", logfile)
			c, err := cmp.Init()
			if err != nil {
				logrus.Fatalf("Can not init http, err %v", err)
			}
			results, ok, err = c.Run(ss)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(3)
			}
			if ok {
				for _, r := range results {
					fmt.Printf("%v:%v: True\n", r.Url, r.Build)
				}
				fmt.Println("------------")
				fmt.Println("Test: passed")
				os.Exit(0)
			} else {
				for _, r := range results {
					fmt.Printf("%v:%v: False\n", r.Url, r.Build)
				}
				fmt.Println("------------")
				fmt.Println("Test: failed")
				os.Exit(1)
			}

		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	logrus.Warnf("Using config file: %v", CfgFile)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&CfgFile, "config", "c", "", "config file (default is $HOME/checkbuild.yaml)")
	rootCmd.Flags().StringVar(&logFile, "logfile", "", "log file (default is Stdout")
	rootCmd.Flags().StringSliceVarP(&urls, "urls", "u", []string{"https://qa.adobeprojectm.com/version", "https://spark.adobe.com/version"}, "Add urls, separated by ','; urls >=2")
	// rootCmd.MarkFlagRequired("urls")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	rootCmd.Flags().BoolVarP(&debug, "debug", "d", false, "debug mode")
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.WarnLevel)
	//logrus.SetReportCaller(true)
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	if CfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(CfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			er(err)
		}

		// Search config in home directory with name "checkbuild" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName("checkbuild")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		logrus.Warnf("Using config file: %v", viper.ConfigFileUsed())
	}
}
