/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/takutakahashi/k2ssm/pkg/k8s"
	"github.com/takutakahashi/k2ssm/pkg/output"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k2ssm",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		o, err := GetParams(cmd, args)
		if err != nil {
			logrus.Fatal(err)
		}
		ctx := context.Background()
		k, err := k8s.NewKubernetes()
		if err != nil {
			logrus.Fatal(err)
		}
		res, err := k.GatherSecrets(ctx, k8s.GatherSecretsOpt{
			Namespace:   "default",
			MatchLabels: map[string]string{},
		})
		if err != nil {
			logrus.Fatal(err)
		}
		if o.sops {
			sops := output.Sops{
				Raw:            res,
				OutputPath:     o.outputPath,
				SopsBinaryPath: "sops",
			}
			if err := sops.Execute(ctx); err != nil {
				logrus.Fatal(err)
			}
		}
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

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.k2ssm.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("config", "c", "./config.yaml", "Configuration file path")
}

type opt struct {
	sops           bool
	secretManager  bool
	externalSecret bool
	outputPath     string
}

func GetParams(cmd *cobra.Command, args []string) (opt, error) {
	return opt{sops: true, secretManager: true, externalSecret: true, outputPath: "./output.yaml"}, nil
}
