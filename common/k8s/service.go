package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-client-go/common"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ReadServiceYaml(filename string) coreV1.Service {
	return GetFromJson(common.ParseYamlToJson(filename))
}

func GetFromJson(jsonByte []byte) coreV1.Service {
	// JSON转struct
	var service coreV1.Service
	err := json.Unmarshal(jsonByte, &service)
	if err != nil {
		panic(err)
	}
	return service
}

func ApplyService(client kubernetes.Clientset, newService coreV1.Service) {
	// 查询k8s是否有该service
	service, err := client.CoreV1().Services("default").Get(context.TODO(), newService.Name, metaV1.GetOptions{})
	if err != nil {
		if !errors.IsNotFound(err) {
			fmt.Println(err)
			return
		}
		// 不存在则创建
		_, err = client.CoreV1().Services("default").Create(context.TODO(), &newService, metaV1.CreateOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
	} else { // 已存在则更新
		//service.Spec.Selector = newService.Spec.Selector
		newService.ResourceVersion = service.ResourceVersion
		newService.Spec.ClusterIP = service.Spec.ClusterIP
		_, err = client.CoreV1().Services("default").Update(context.TODO(), &newService, metaV1.UpdateOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Printf("apply service %s success!", service.Name)
}
