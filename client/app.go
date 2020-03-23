package client

import (
	"github.com/spf13/cobra"

	"github.com/qingstor/openvpn-warder/check"
	"github.com/qingstor/openvpn-warder/config"
	"github.com/qingstor/openvpn-warder/constants"
)

var (
	// client
	flagConfig string

	// get user
	getUserName string
)

// Client respresent openvpn-warder client app.
var Client = &cobra.Command{
	Use:   "client",
	Short: "openvpn-warder CLI",
	Long:  "The API CLI of openvpn-warder server",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewClient()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
	},
}

var getUser = &cobra.Command{
	Use:   "get-user",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewClient()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		GetUser(c, getUserName)
	},
}

// InitClient will init client app.
func InitClient() {
	Client.SilenceErrors = true
	Client.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	Client.PersistentFlags().StringVarP(
		&flagConfig, "config", "c", constants.Name, "Specify client config file path",
	)

    initGetUser()
	Client.AddCommand(getUser)
}

func initGetUser() {
	getUser.SilenceErrors = true
	getUser.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	getUser.Flags().StringVarP(
		&getUserName, "name", "", constants.Name, "Specify user name you want to get",
	)
}
