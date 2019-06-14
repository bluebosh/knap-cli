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

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all knap appengines",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := clientcmd.BuildConfigFromFlags("", "/Users/jordan/.bluemix/plugins/container-service/clusters/knative_pipeline/kube-config-dal10-knative_pipeline.yml")
		if err != nil {
			glog.Fatalf("Error building kubeconfig: %v", err)
		}

		knapClient, err := knapclientset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building knap clientset: %v", err)
		}

		appLst, err := knapClient.KnapV1alpha1().Appengines("default").List(metav1.ListOptions{})
		color.Cyan("%-30s%-20s%-20s%-20s\n", "Application Name", "Status", "Instance", "Domain")
		for i := 0; i < len(appLst.Items); i++ {
			app := appLst.Items[i]
			fmt.Printf("%-30s%-20s%-20s%-20s\n", app.Spec.AppName + "-appengine", app.Status.Status, "1/1", app.Status.Domain + "\n")
		}
		fmt.Println("There are", color.CyanString("%v",len(appLst.Items)), "application engine(s)")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
