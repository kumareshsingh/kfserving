/*
Copyright 2019 kubeflow.org.

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

package v1alpha2

import (
	"context"
	"github.com/kubeflow/kfserving/pkg/constants"
	"k8s.io/api/core/v1"
	"os"
	"path/filepath"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var cfg *rest.Config
var c client.Client

const (
	DefaultTensorflowRuntimeVersion     = "latest"
	DefaultTensorflowRuntimeVersionGPU  = "latest-gpu"
	DefaultSKLearnRuntimeVersion        = "0.1.0"
	DefaultPyTorchRuntimeVersion        = "0.1.0"
	DefaultPyTorchRuntimeVersionGPU     = "0.1.0-gpu"
	DefaultXGBoostRuntimeVersion        = "0.1.0"
	DefaultTensorRTRuntimeVersion       = "19.05-py3"
	DefaultONNXRuntimeVersion           = "v0.5.0"
	DefaultAlibiExplainerRuntimeVersion = "0.2.3"
)

func TestMain(m *testing.M) {
	t := &envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "..", "..", "..", "config", "default", "crds", "base")},
	}

	err := SchemeBuilder.AddToScheme(scheme.Scheme)

	if err != nil {
		klog.Fatal(err)
	}

	if cfg, err = t.Start(); err != nil {
		klog.Fatal(err)
	}

	if c, err = client.New(cfg, client.Options{Scheme: scheme.Scheme}); err != nil {
		klog.Fatal(err)
	}

	// Create configmap
	configs := map[string]string{
		"predictors": `{
			"tensorflow" : {
				"image" : "tensorflow/serving",
				"defaultImageVersion": "latest",
				"defaultGPUImageVersion": "latest-gpu",
				"allowedImageVersions": [
				   "latest",
				   "latest-gpu"
				]
			},
			"sklearn" : {
				"image" : "kfserving/sklearnserver",
				"defaultImageVersion": "0.1.0",
				"allowedImageVersions": [
				   "latest",
				   "0.1.0"
				]
			},
			"xgboost" : {
				"image" : "kfserving/xgbserver",
				"defaultImageVersion": "0.1.0",
				"allowedImageVersions": [
				   "latest",
				   "0.1.0"
				]
			},
			"pytorch" : {
				"image" : "kfserving/pytorchserver",
				"defaultImageVersion": "0.1.0",
                "defaultGPUImageVersion": "0.1.0-gpu",
				"allowedImageVersions": [
				   "latest",
				   "0.1.0",
                   "0.1.0-gpu"
				]
			},
			"onnx" : {
				"image" : "onnxruntime/server",
				"defaultImageVersion": "v0.5.0",
				"allowedImageVersions": [
				   "latest",
				   "v0.5.0"
				]
			},
			"tensorrt" : {
				"image" : "nvcr.io/nvidia/tensorrtserver",
				"defaultImageVersion": "19.05-py3",
				"allowedImageVersions": [
				   "19.05-py3"
				]
			}
		}`,
		"explainers": `{
			"alibi" : {
				"image" : "docker.io/seldonio/alibiexplainer",
				"defaultImageVersion": "0.2.3",
				"allowedImageVersions": [
				   "0.2.3"
				]
			}
        }`,
	}
	var configMap = &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      constants.InferenceServiceConfigMapName,
			Namespace: constants.KFServingNamespace,
		},
		Data: configs,
	}
	if err := c.Create(context.TODO(), configMap); err != nil {
		klog.Fatal(err)
	}
	defer c.Delete(context.TODO(), configMap)

	code := m.Run()
	t.Stop()
	os.Exit(code)
}
