/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	azureinfo "github.com/Rbsn-joses/create-pbi/azureInfo"
	"github.com/Rbsn-joses/create-pbi/customLogger"
	"github.com/spf13/cobra"
)

// azureDevopsInfoCmd represents the injectPbis command
var azureDevopsInfoCmd = &cobra.Command{
	Use:   "azureDevopsInfo",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := customLogger.InitLogger(Verbose)

		azureinfo.StartInjection(UsernameDevops, PAT, Project, OrganizationURL, ExcelFile, logger)
	},
}

func init() {
	rootCmd.AddCommand(azureDevopsInfoCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&PAT, "personalAccessToken", "t", "", "Token de acesso do azureDevops")
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "azureDevops project")
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&UsernameDevops, "usernameDevops", "u", "", "azureDevops username")
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&OrganizationURL, "organizationURL", "o", "", "azureDevops organization url")
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&ExcelFile, "excelFile", "f", "", "path excelfile")
	azureDevopsInfoCmd.PersistentFlags().StringVarP(&Verbose, "verbose", "v", "info", "saida com mais informações para debug")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// azureDevopsInfoCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
