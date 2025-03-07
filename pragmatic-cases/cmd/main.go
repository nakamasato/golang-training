package main

import (
	"fmt"
	"github.com/spf13/cobra"
)


var cmd = &cobra.Command{
	Use:   "check",
	Short: "check",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("PARENT")
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		fmt.Println("PARENT POST")
		for _, f := range deferFuncs {
			defer f()
		}
	},
}

var chidCmd = &cobra.Command{
	Use:   "child",
	Short: "child",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CHILD")
	},
}

func main() {
	cmd.AddCommand(chidCmd)
	cmd.Execute()
	fmt.Println("main")
}

var deferFuncs []func()

func init() {
	fmt.Println("init")
	deferFuncs = append(deferFuncs, func() {
		fmt.Println("defer")
	})
}

