package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"encoding/json"

	"github.com/spf13/cobra"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/validation"
	"k8s.io/apimachinery/pkg/conversion"
)

var crdPath string
var debug bool

func hasAnyStatusEnabled(crd *apiextensions.CustomResourceDefinitionSpec) bool {
	if hasStatusEnabled(crd.Subresources) {
		return true
	}
	for _, v := range crd.Versions {
		if hasStatusEnabled(v.Subresources) {
			return true
		}
	}
	return false
}

// hasStatusEnabled returns true if given CRD Subresources has non-nil Status set.
func hasStatusEnabled(subresources *apiextensions.CustomResourceSubresources) bool {
	if subresources != nil && subresources.Status != nil {
		return true
	}
	return false
}

var RootCmd = &cobra.Command{
	Use:   "validate",
	Short: "for validation of crd",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		jsonFile, err := ioutil.ReadFile(crdPath)
		if err != nil {
			panic(err)
		}

		var obj v1beta1.CustomResourceDefinition
		var apiObj apiextensions.CustomResourceDefinition
		var s conversion.Scope
		json.Unmarshal([]byte(jsonFile), &obj)
		if debug {
			fmt.Printf("%+v\n\n", obj)
		}

		v1beta1.Convert_v1beta1_CustomResourceDefinition_To_apiextensions_CustomResourceDefinition(&obj, &apiObj, s)
		if debug {
			fmt.Printf("%+v\n\n", apiObj.Spec.Validation.OpenAPIV3Schema)
		}

		// version v0.0.0-20190918201827-3de75813f604
		errList := validation.ValidateCustomResourceDefinition(&apiObj)
		for _, e := range errList {
			fmt.Printf("%+v\n", e)
		}

		// version v0.18.2
		// requestGV := schema.GroupVersion{
		// 	Group:   apiObj.Spec.Group,
		// 	Version: apiObj.Spec.Version,
		// }
		// errList := validation.ValidateCustomResourceDefinition(&apiObj, requestGV)

	},
}

func init() {
	RootCmd.PersistentFlags().StringVar(&crdPath, "def", "", "path to crd definition")
	RootCmd.PersistentFlags().BoolVar(&debug, "d", false, "Help message for toggle")
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
