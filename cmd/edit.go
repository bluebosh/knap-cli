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
	"strconv"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit a knap appengine",
	Long: `Edit a knap appengine`,
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
			fmt.Println("Error getting application engine", color.CyanString(args[0]), err)
		} else {
			size, err:= strconv.ParseInt(cmd.Flag("size").Value.String(),10,32)
			size32 := int32(size)
			if err != nil {
				//glog.Fatalf("Error creating application engine: %s", args[0])
				fmt.Println("Error parsing size parameter", err)
			}

			app.Spec.GitRepo = cmd.Flag("gitrepo").Value.String()
			app.Spec.GitRevision = cmd.Flag("gitrevision").Value.String()
			//app.Spec.GitWatch = cmd.Flag("getwatch").Value.String()
			app.Spec.Size = size32
			app.Spec.PipelineTemplate = cmd.Flag("template").Value.String()

			app, err = knapClient.KnapV1alpha1().Appengines("default").Update(app)
			if err != nil {
				//glog.Fatalf("Error creating application engine: %s", args[0])
				fmt.Println("Error updating application engine", color.CyanString(args[0]), err)
			} else {
				fmt.Println("Application engine", color.CyanString(args[0]), "is updated successfully")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	editCmd.Flags().StringP("resourcetype","e","", "The resource type the appengine [Not implemented]")
	editCmd.Flags().StringP("gitrepo","r","", "The git repo of the appengine")
	editCmd.Flags().StringP("gitrevision","v","", "The git revision of the appengine")
	editCmd.Flags().StringP("template","t","", "The template of the appengine")
	editCmd.Flags().Int32P("size","s",1, "The size of the appengine")
	editCmd.Flags().BoolP("watch", "a", false, "The auto trigger of the appengine [Not implemented]")
}
