package harborv1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const GroupName = "harbor.f110.dev"

var (
	GroupVersion       = metav1.GroupVersion{Group: GroupName, Version: "v1alpha1"}
	SchemeBuilder      = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme        = SchemeBuilder.AddToScheme
	SchemaGroupVersion = schema.GroupVersion{Group: "harbor.f110.dev", Version: "v1alpha1"}
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemaGroupVersion,
		&HarborProject{},
		&HarborProjectList{},
		&HarborRobotAccount{},
		&HarborRobotAccountList{},
	)
	metav1.AddToGroupVersion(scheme, SchemaGroupVersion)
	return nil
}