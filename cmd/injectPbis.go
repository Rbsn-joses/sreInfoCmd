/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/Rbsn-joses/create-pbi/customLogger"
	"github.com/Rbsn-joses/create-pbi/task"
	"github.com/spf13/cobra"
)

var Project string
var OrganizationURL string
var PAT string
var ExcelFile string
var UsernameDevops string
var Verbose string

// injectPbisCmd represents the injectPbis command
var injectPbisCmd = &cobra.Command{
	Use:   "injectPBI",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := customLogger.InitLogger(Verbose)

		task.StartInjection(UsernameDevops, PAT, Project, OrganizationURL, ExcelFile, logger)
	},
}

func init() {
	rootCmd.AddCommand(injectPbisCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	injectPbisCmd.PersistentFlags().StringVarP(&PAT, "personalAccessToken", "t", "", "Token de acesso do azureDevops")
	injectPbisCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "azureDevops project")
	injectPbisCmd.PersistentFlags().StringVarP(&UsernameDevops, "usernameDevops", "u", "", "azureDevops username")
	injectPbisCmd.PersistentFlags().StringVarP(&OrganizationURL, "organizationURL", "o", "", "azureDevops organization url")
	injectPbisCmd.PersistentFlags().StringVarP(&ExcelFile, "excelFile", "f", "", "path excelfile")
	injectPbisCmd.PersistentFlags().StringVarP(&Verbose, "verbose", "v", "info", "saida com mais informações para debug")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// injectPbisCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
