package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"

	"github.com/qingstor/openvpn-warder/check"
	"github.com/qingstor/openvpn-warder/client"
	"github.com/qingstor/openvpn-warder/config"
	"github.com/qingstor/openvpn-warder/constants"
	"github.com/qingstor/openvpn-warder/handlers"
	"github.com/qingstor/openvpn-warder/log"
)

var (
	flagVersion bool
	flagConfig  string
)

var application = &cobra.Command{
	Use:   constants.Name,
	Short: "openvpn-warder",
	Long:  "The API Server of openvpn-warder",
	Run: func(cmd *cobra.Command, args []string) {
		if flagVersion {
			fmt.Printf("%s version %s\n", constants.Name, constants.Version)
			return
		}

		c := config.NewWarderServer()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		config.DBPath = *c.DBPath

        addr := fmt.Sprintf("0.0.0.0:%d", *c.Port)
		app := gin.Default()
		app.POST("/users/get", handlers.HandleGetUser)
		err = app.Run(addr)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
		}
		log.Logger.Info(
			"Openvpn warder started\n",
			fmt.Sprintf("Listenning at %s\n", addr),
			fmt.Sprintf("User update cycle is %d days\n", c.UpdateCycle),
		)
	},
}

func init() {
	application.SilenceErrors = true
	application.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	application.Flags().BoolVarP(
		&flagVersion, "version", "v", false, "Show version",
	)
	application.Flags().StringVarP(
		&flagConfig, "config", "c", constants.Name, "Specify config file path",
	)

	client.InitClient()
	application.AddCommand(client.Client)
}

func main() {
	check.ErrorForExit(constants.Name, application.Execute())
}
