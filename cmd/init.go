package cmd

func init() {
	RootCmd.Flags().BoolP("init", "i", false, "init config")
	// RootCmd.AddCommand(initCommand)
}

// var initCommand = &cobra.Command{
// 	Use:   "init",
// 	Short: "init config",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("init config")
// 	},
// }
