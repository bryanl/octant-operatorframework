module github.com/bryanl/octant-operatorframework

go 1.13

require (
	github.com/golang/mock v1.3.1
	github.com/golang/protobuf v1.3.2
	github.com/google/uuid v1.1.1
	github.com/operator-framework/operator-lifecycle-manager v0.0.0-20200213132121-99643127d862
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/vmware-tanzu/octant v0.9.2-0.20191221004819-6a2b14b497a1
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	google.golang.org/grpc v1.27.0
	k8s.io/api v0.17.1
	k8s.io/apimachinery v0.17.1
)

replace (
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200217161739-c99157bc6492
	github.com/vmware-tanzu/octant => /Users/bryan/Development/projects/octant
	k8s.io/api => k8s.io/api v0.16.7
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.16.7
	k8s.io/apimachinery => k8s.io/apimachinery v0.16.7
	k8s.io/apiserver => k8s.io/apiserver v0.16.7
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.16.7
	k8s.io/client-go => k8s.io/client-go v0.16.7
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.16.7
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.16.7
	k8s.io/code-generator => k8s.io/code-generator v0.16.7
	k8s.io/component-base => k8s.io/component-base v0.16.7
	k8s.io/cri-api => k8s.io/cri-api v0.16.7
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.16.7
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.16.7
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.16.7
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.16.7
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.16.7
	k8s.io/kubectl => k8s.io/kubectl v0.16.7
	k8s.io/kubelet => k8s.io/kubelet v0.16.7
	k8s.io/kubernetes => k8s.io/kubernetes v1.16.7
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.16.7
	k8s.io/metrics => k8s.io/metrics v0.16.7
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.16.7
)
