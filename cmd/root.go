/*
Copyright © 2023 Tomasz Sobota <tomasz@sobota.cc>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"csvtoredis/pkg"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "csvtoredis",
	Short: "Upload a CSV file to Redis",
	Long: `Upload a CSV file to Redis. The tool will create a key for each column and row. 
				The key name will be in the format <prefix><column_name>-<row_number>.`,
	Run: func(cmd *cobra.Command, args []string) {

		csv_source := pkg.NewCSVSource()
		err := csv_source.ReadCSV(viper.GetString("csv"))
		if err != nil {
			log.Println("Error reading CSV file")
			log.Println(err)
			return
		}

		just_text := viper.GetBool("just-text")
		if just_text {
			csv_source.Print(viper.GetString("redis-key-prefix"))
			return
		}

		redis_target := pkg.NewRedisTarget()
		redis_target.Host = viper.GetString("redis-host")
		redis_target.Port = viper.GetString("redis-port")
		redis_target.Password = viper.GetString("redis-password")

		err = redis_target.WriteToRedis(csv_source, viper.GetString("redis-key-prefix"))
		if err != nil {
			log.Println("Error writing to Redis")
			log.Println(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.csvtoredis.yaml)")
	rootCmd.PersistentFlags().StringP("csv", "s", "", "CSV file to read from")
	rootCmd.PersistentFlags().StringP("redis-host", "r", "localhost", "Redis host")
	rootCmd.PersistentFlags().StringP("redis-port", "p", "6379", "Redis port")
	rootCmd.PersistentFlags().StringP("redis-password", "w", "", "Redis password")
	rootCmd.PersistentFlags().StringP("redis-key-prefix", "k", "", "Redis key prefix")
	rootCmd.PersistentFlags().Bool("just-text", false, "Plain text output (dry-run)")

	viper.BindPFlag("csv", rootCmd.PersistentFlags().Lookup("csv"))
	viper.BindPFlag("redis-host", rootCmd.PersistentFlags().Lookup("redis-host"))
	viper.BindPFlag("redis-port", rootCmd.PersistentFlags().Lookup("redis-port"))
	viper.BindPFlag("redis-password", rootCmd.PersistentFlags().Lookup("redis-password"))
	viper.BindPFlag("redis-key-prefix", rootCmd.PersistentFlags().Lookup("redis-key-prefix"))
	viper.BindPFlag("just-text", rootCmd.PersistentFlags().Lookup("just-text"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".csvtoredis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".csvtoredis")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
