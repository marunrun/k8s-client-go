package istio

import (
	"context"
	"encoding/json"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	"k8s-client-go/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func GetVirtualServicesByName(client versioned.Clientset, serviceName string, namespace string) (*v1alpha3.VirtualService, error) {
	if namespace == "" {
		namespace = common.DefaultNamespace
	}
	return client.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), serviceName, metav1.GetOptions{})
}

// GetVirtualServicesFromYamlFile 将yaml文件转换成VirtualServices
func GetVirtualServicesFromYamlFile(fileName string) v1alpha3.VirtualService {
	return GetVirtualServicesFromJson(common.ParseYamlToJson(fileName))
}

// GetVirtualServicesFromJson  将json转换成VirtualServices
func GetVirtualServicesFromJson(jsonByte []byte) v1alpha3.VirtualService {
	var service v1alpha3.VirtualService
	err := json.Unmarshal(jsonByte, &service)
	if err != nil {
		panic(err.Error())
	}
	return service
}

// CreateVirtualServices 新建 VirtualServices
func CreateVirtualServices(client versioned.Clientset, service v1alpha3.VirtualService, namespace string) (*v1alpha3.VirtualService, error) {
	return client.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), &service, metav1.CreateOptions{})
}

