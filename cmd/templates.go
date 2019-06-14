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
	tektoncdentset "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"github.com/fatih/color"
)

// templatesCmd represents the templates command
var templatesCmd = &cobra.Command{
	Use:   "templates",
	Short: "List all knap templates",
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

		tektoncdClient, err := tektoncdentset.NewForConfig(cfg)
		if err != nil {
			glog.Fatalf("Error building knap clientset: %v", err)
		}

		pipelines, err := tektoncdClient.TektonV1alpha1().Pipelines("default").List(metav1.ListOptions{})
		color.Cyan("%-40s%-80s\n", "Template Name", "Template Flow")
		for i := 0; i < len(pipelines.Items); i++ {
			pipeline := pipelines.Items[i]
			taskFlow := ""
			for i := 0; i < len(pipeline.Spec.Tasks); i++ {
				task := pipeline.Spec.Tasks[i]
				if taskFlow == "" {
					taskFlow = task.Name
				} else {
					taskFlow = taskFlow + " -> " + task.Name
				}
			}
			fmt.Printf("%-40s%-80s\n", pipeline.Name, taskFlow)
		}
	},
}

func init() {
	rootCmd.AddCommand(templatesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// templatesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// templatesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
