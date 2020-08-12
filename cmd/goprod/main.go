package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/andrewhowdencom/gopro/internal/hotplug"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "goprod",
	Short: "Daemon that manages the connection to the GoPro",
	Run: func(cmd *cobra.Command, args []string) {
		h, e := hotplug.New()
		if e != nil {
			log.Fatal(e.Error())
		}

		c, e := h.Listen(context.Background())
		if e != nil {
			log.Fatalf("unable to listen for events: %s", e.Error())
		}

		for event := range c {
			switch event.Type {
			case hotplug.Connected:
				fmt.Printf("Yeah! Found GoPro with ID %s\n", event.Entity.ID)
			case hotplug.Disconnected:
				fmt.Printf("Boo! Camera with id %s went away\n", event.Entity.ID)
			}
		}
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is /etc/goprod.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {

		// Todo: Make these paths based on XDG configuration
		viper.AddConfigPath("/etc")
		viper.SetConfigName(".goprod")
	}

	viper.SetEnvPrefix("GOPROD")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
