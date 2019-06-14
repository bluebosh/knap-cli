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
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [appengine name]",
	Short: "Get a knap appengine detail",
	Long: `Get a knap appengine detail`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Println("Error getting application engine", color.CyanString(args[0]))
		} else {
			fmt.Println(color.CyanString("%-30s","Application Name:"), app.Spec.AppName + "-appengine")
			fmt.Println(color.CyanString("%-30s","Application Version:"), app.Generation)
			fmt.Println(color.CyanString("%-30s","Application Git Repo:"), app.Spec.GitRepo)
			fmt.Println(color.CyanString("%-30s","Application Git Revision:"), app.Spec.GitRevision)
			fmt.Println(color.CyanString("%-30s","Application Template:"), app.Spec.PipelineTemplate)
			fmt.Println(color.CyanString("%-30s","Application Status:"), app.Status.Status)
			fmt.Println(color.CyanString("%-30s","Application Instance:"), "1")
			fmt.Println(color.CyanString("%-30s","Application Domain:"),"https://" + app.Status.Domain)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
