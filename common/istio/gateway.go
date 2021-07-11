package istio

import (
	"context"
	"encoding/json"
	"istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/client-go/pkg/clientset/versioned"
	"k8s-client-go/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)





func GetGatewayFromYamlFile(fileName string) v1alpha3.VirtualService {
	return GetGatewayFromJson(common.ParseYamlToJson(fileName))
}

func GetGatewayFromJson(jsonByte []byte) v1alpha3.VirtualService {
	var service v1alpha3.VirtualService
	err := json.Unmarshal(jsonByte, &service)
	if err != nil {
		panic(err.Error())
	}
	return service
}

func CreateGateway(client versioned.Clientset, service v1alpha3.VirtualService, namespace string) (*v1alpha3.VirtualService, error) {
	return client.NetworkingV1alpha3().VirtualServices(namespace).Create(context.TODO(), &service, metav1.CreateOptions{})
}



func GetGatewayByName(client versioned.Clientset, name string, namespace string) (*v1alpha3.Gateway, error) {
	if namespace == "" {
		namespace = common.DefaultNamespace
	}

	return  client.NetworkingV1alpha3().Gateways(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}


