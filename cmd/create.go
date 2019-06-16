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
	knapv1 "github.com/bluebosh/knap/pkg/apis/knap/v1alpha1"
	"github.com/golang/glog"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc" // from https://github.com/kubernetes/client-go/issues/345
	"github.com/fatih/color"
	"strconv"
)

var foo string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [appengine name]",
	Short: "Create a new knap appengine",
	Long: `Create a new knap appengine`,
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


		size, err:= strconv.ParseInt(cmd.Flag("size").Value.String(),10,32)
		size32 := int32(size)
		if err != nil {
			//glog.Fatalf("Error creating application engine: %s", args[0])
			fmt.Println("Error parsing size parameter", err)
		}

		app := &knapv1.Appengine{
			ObjectMeta: metav1.ObjectMeta{
				Name:      args[0] + "-appengine",
				Namespace: "default",
			},
			Spec:
			knapv1.AppengineSpec{
				AppName: args[0],
				GitRepo: cmd.Flag("gitrepo").Value.String(),
				GitRevision: cmd.Flag("gitrevision").Value.String(),
				// GitWatch: cmd.Flag("gitWatch").Value.String(),
				Size: size32,
				PipelineTemplate: cmd.Flag("template").Value.String(),
			},
		}
		_, err = knapClient.KnapV1alpha1().Appengines("default").Create(app)

		if err != nil {
			//glog.Fatalf("Error creating application engine: %s", args[0])
			fmt.Println("Error creating application engine", color.CyanString(args[0]), err)
		} else {
			fmt.Println("Application engine", color.CyanString(args[0]), "is created successfully")
		}
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().StringP("gitrepo","r","", "The git repo of the appengine")
	createCmd.Flags().StringP("gitrevision","v","", "The git revision of the appengine")
	createCmd.Flags().StringP("template","t","", "The template of the appengine")
	createCmd.Flags().Int32P("size","s",1, "The size of the appengine")
	createCmd.Flags().BoolP("autotrigger", "a", false, "The auto trigger of the appengine")
}
