/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var ErrorColor string = "\033[38;5;9m"
var NoColor string = "\033[0m"
var SuccessColor string = "\033[38;5;10m"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "homelab-cli",
	Short: "A CLI tool to manage your homelab",
	Long:  `homelab-cli aims to be a simple solution to create new nixos vms on proxmox. Together with https://github.com/NullNUMMER24/nix-server-conf-manager managing a homelab with nixos should be a breeze.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.homelab-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
