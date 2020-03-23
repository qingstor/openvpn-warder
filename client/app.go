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

	// create user
	createUserName     string
	createUserPassword string
	createUserAdmin    bool
	createUserIgnore   bool

	// delete user
	deleteUserName string

	// reset user
	resetUserName        string
	resetUserNewName     string
	resetUserNewPassword string
	resetUserNewAdmin    bool
	resetUserNewIgnore   bool
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

var createUser = &cobra.Command{
	Use:   "create-user",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewClient()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		CreateUser(
			c,
			createUserName,
			createUserPassword,
			createUserAdmin,
			createUserIgnore,
		)
	},
}

var deleteUser = &cobra.Command{
	Use:   "delete-user",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewClient()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		DeleteUser(
			c,
			deleteUserName,
		)
	},
}

var resetUser = &cobra.Command{
	Use:   "reset-user",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		c := config.NewClient()
		err := c.LoadFromFilePath(flagConfig)
		if err != nil {
			check.ErrorForExit(constants.Name, err)
			return
		}
		ResetUser(
			c,
			resetUserName,
			resetUserNewName,
			resetUserNewPassword,
			resetUserNewAdmin,
			resetUserNewIgnore,
		)
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
	initCreateUser()
	Client.AddCommand(createUser)
	initDeleteUser()
	Client.AddCommand(deleteUser)
	initResetDeleteUser()
	Client.AddCommand(resetUser)
}

func initGetUser() {
	getUser.SilenceErrors = true
	getUser.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	getUser.Flags().StringVarP(
		&getUserName, "name", "", "", "Specify user name you want to get",
	)
}

func initCreateUser() {
	createUser.SilenceErrors = true
	createUser.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	createUser.Flags().StringVarP(
		&createUserName, "name", "", "", "Specify user name you want to create",
	)
	createUser.Flags().StringVarP(
		&createUserPassword, "password", "", "", "Specify user password you want to create",
	)
	createUser.Flags().BoolVarP(
		&createUserAdmin, "admin", "", false, "Specify user admin you want to create",
	)
	createUser.Flags().BoolVarP(
		&createUserIgnore, "ignore", "", false, "Specify user ignore you want to create",
	)
}

func initDeleteUser() {
	deleteUser.SilenceErrors = true
	deleteUser.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	deleteUser.Flags().StringVarP(
		&deleteUserName, "name", "", "", "Specify user name you want to delete",
	)
}

func initResetDeleteUser() {
	resetUser.SilenceErrors = true
	resetUser.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		check.ErrorForExit(constants.Name, err)
		return nil
	})

	resetUser.Flags().StringVarP(
		&resetUserName, "name", "", "", "Specify user name you want to reset",
	)
	resetUser.Flags().StringVarP(
		&resetUserNewName, "new-name", "", "", "Specify user new name you want to reset",
	)
	resetUser.Flags().StringVarP(
		&resetUserNewPassword, "new-password", "", "", "Specify user new password you want to reset",
	)
	resetUser.Flags().BoolVarP(
		&resetUserNewAdmin, "new-admin", "", false, "Specify user new admin you want to reset",
	)
	resetUser.Flags().BoolVarP(
		&resetUserNewIgnore, "new-ignore", "", false, "Specify user new ignore you want to reset",
	)
}
