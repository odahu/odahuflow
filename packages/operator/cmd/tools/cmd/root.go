package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/odahu/odahu-flow/packages/operator/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"runtime/pprof"
)

var (
	cfgFile string
	cfg config.ToolsConfig
	profilingFile *os.File
	apiURL string
)

var rootCmd = &cobra.Command{
	Use:   "odahu-tools",
	Short: "odahu-tools is a simple command line tool that provides API to the set of ODAHU platform features",
	Long: `odahu-tools provides API to execute the same logic that is used by the ODAHU platform in the cluster`,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			os.Exit(0)
		}
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			var err error
			profilingFile, err = os.Create(cpuprofile)
			if err != nil {
				zap.S().Fatal("could not create CPU profile: ", err)
			}
			if err := pprof.StartCPUProfile(profilingFile); err != nil {
				zap.S().Fatal("could not start CPU profile: ", err)
			}
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if cpuprofile != "" {
			pprof.StopCPUProfile()
			if profilingFile != nil {
				if err := profilingFile.Close(); err != nil {
					zap.S().Fatal(err)
				}
			}
		}

	},
}

func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal("Unable to initialize logger")
	}
	defer func() {
		_ = logger.Sync()
	}()

	zap.ReplaceGlobals(logger)
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile == "" {
		cfgFile = os.Getenv("ODAHU_TOOLS_CONFIG")
	}
	if cfgFile != "" {
		// Use config file from the flag.
		zap.S().Infof("Config path: %s", cfgFile)
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		zap.S().Info("Config path is not configured explicitly. Looking for a config in home directory")
		home, err := homedir.Dir()
		if err != nil {
			zap.S().Errorw("Unable to get home directory", zap.Error(err))
			panic("Unable to get home directory")
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".odahu-tools")
	}

	if err := viper.ReadInConfig(); err != nil {
		zap.S().Errorw("Can't read config", zap.Error(err))
		panic("Can't read config")
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		zap.S().Errorw("Unable to unmarshall config", zap.Error(err))
		panic(err)
	}

}

func init() {
	initLogger()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(
		&cfgFile, "config", "", "config file (default is $HOME/.odahu-tools.yaml)",
	)
	rootCmd.PersistentFlags().StringVar(
		&cpuprofile, "cpuprofile", "",
		"if specified as TARGET cpu profiling results will be saved to TARGET file",
	)
	rootCmd.PersistentFlags().StringVar(
		&apiURL, "api", "", "API server base URL (schema://host:port)",
	)
	_ = viper.BindPFlag("auth.apiUrl", rootCmd.PersistentFlags().Lookup("api"))
	// These env variables are compatible with odahu-flow-sdk package
	_ = viper.BindEnv("auth.apiUrl", "API_URL")
	_ = viper.BindEnv("auth.apiToken", "API_TOKEN")
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("Exit code 1: %s\n", err)
		os.Exit(1)
	}
}