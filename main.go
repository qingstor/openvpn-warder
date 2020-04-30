package main

import (
	"fmt"
	"os"

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

		// Client
		userName := os.Getenv("username")
		userPassword := os.Getenv("password")
		if userName != "" && userPassword != "" {
			clientConfig := config.NewClient()
			err := clientConfig.LoadFromFilePath(config.DefaultClientConfigPath)
			if err != nil {
				check.ErrorForExit(constants.Name, err)
			}
			client.VerifyUser(clientConfig, userName, userPassword)
		}

		c := config.NewWarderServer()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		config.DBPath = *c.DBPath
		err = log.InitLogger(*c.LogPath)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
		}

		err = handlers.CycleChangePassword(c.Mail, *c.UpdateCycle, *c.VPNName)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
		}

		addr := fmt.Sprintf("0.0.0.0:%d", *c.Port)
		app := gin.Default()
		app.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
		app.POST("/users/get", handlers.HandleGetUser)
		app.POST("/users/create", handlers.HandleCreateUser)
		app.POST("/users/delete", handlers.HandleDeleteUser)
		app.POST("/users/reset", handlers.HandleResetUser)
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
