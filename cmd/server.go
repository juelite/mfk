package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"mfk/logic"
)

var newCmd = &cobra.Command{
	Use:     "new",
	Example: "mfk new -n blog",
	Short:   "创建一个新的项目",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : ", err)
			}
		}()
		logic.New()
	},
}

var runCmd = &cobra.Command{
	Use:     "run",
	Example: "mfk run",
	Short:   "运行项目",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : ", err)
			}
		}()
		logic.Run()
	},
}

var packCmd = &cobra.Command{
	Use:     "pack",
	Example: "mfk pack",
	Short:   "项目打包",
	Run: func(cmd *cobra.Command, args []string) {
		defer func() {
			if err := recover(); err != nil {
				log.Println("Recover error : ", err)
			}
		}()
		logic.Pack()
	},
}

func init() {
	newCmd.Flags().StringVarP(&logic.NewProject, "name", "n", "", "项目名称")
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(packCmd)
}
