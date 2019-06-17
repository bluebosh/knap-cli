/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	knapclientset "github.com/bluebosh/knap/pkg/client/clientset/versioned"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"github.com/fatih/color"
	"os"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Get a knap appengine logs",
	Long: `Get a knap appengine logs`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		kubeconfig = rootCmd.Flag("kubeconfig").Value.String()
		if kubeconfig == "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			glog.Fatalf("Error building kubeconfig: %v", err)
		}

		knapClient, err := knapclientset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building knap clientset: %v", err)
		}

		app, err := knapClient.KnapV1alpha1().Appengines("default").Get(args[0], metav1.GetOptions{})

		if err != nil {
			fmt.Println("Error getting application engine", color.CyanString(args[0]), err)
		} else {
			fmt.Println("Retrieving application", color.CyanString(app.Spec.AppName), "logs...")
			fmt.Println(app.Spec)
			fmt.Println(app.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(logsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
