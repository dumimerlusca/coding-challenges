/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	cliapp "cut/internal/cliapp"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

const TAB_DELIMITER = "	"

const DEFAULT_DELIMITER string = TAB_DELIMITER

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cut",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// memProf, _ := os.Create("mem.pprof")
		// defer pprof.WriteHeapProfile(memProf)
		// defer memProf.Close()

		fields, err := cmd.Flags().GetString("fields")
		if err != nil {
			panic(err)
		}

		delimiter, err := cmd.Flags().GetString("delimiter")
		if err != nil {
			panic(err)
		}

		info, err := os.Stdin.Stat()
		if err != nil {
			panic(fmt.Errorf("failed to read os.Stdin.Stat %w", err))
		}

		useStdin := (info.Mode() & os.ModeCharDevice) == 0
		var inputReader io.Reader

		if useStdin {
			inputReader = os.Stdin
		} else {
			if len(args) < 2 {
				fmt.Println("filename not provided")
				os.Exit(0)
			}

			filename := args[1]

			file, err := os.Open(filename)
			if err != nil {
				fmt.Println("failed to open file")
				os.Exit(0)
			}

			inputReader = file

		}

		cfg := cliapp.Config{Delimeter: delimiter, Select: cliapp.SelectOption{OptionType: cliapp.SelectOptionFields, Value: fields}}

		a := cliapp.NewApp(cfg, inputReader)

		output, err := a.Run()
		if err != nil {
			fmt.Printf("something went wrong: %v \n", err.Error())
			os.Exit(0)
		}

		fmt.Print(string(output))
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cut.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringP("fields", "f", "", "Fields")
	rootCmd.Flags().StringP("delimiter", "d", DEFAULT_DELIMITER, "Delimeter")
}
