package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"mfk/logic"
)

var serverCmd = &cobra.Command{
	Use:   		"new",
	Example: 	"mfk new -n blog  创建一个名为[blog]的项目",
	Short: "创建一个新的项目",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : ", err)
			}
		}()
		logic.Run()
	},
}

func init() {
	serverCmd.Flags().StringVarP(&logic.NewProject, "name", "n", "", "项目名称")
	rootCmd.AddCommand(serverCmd)
}