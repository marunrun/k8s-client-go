package main

import (
	"flag"
	"fmt"
	"k8s-client-go/common"
	"k8s-client-go/common/istio"
)

func main() {

	//resource := flag.String("resource", "", "需要操作的资源,[deployment,pod,service]")
	//action := flag.String("action", "", "动作,[]")

	flag.Parse()

	// 初始化k8s客户端
	//clientset := common.GetK8sClient()

	//region 创建虚拟服务
	vservice := istio.GetVirtualServicesFromYamlFile("../yamls/hyperf-vService.yaml")
	istioClient := common.GetIstioClient()
	services, err := istio.CreateVirtualServices(*istioClient, vservice, vservice.GetNamespace())
	if err != nil {
		panic(err)
	}
	fmt.Printf("success %s", services.Name)
	//endregion

	//
	////region service
	//service := k8s.ReadServiceYaml("../yamls/hyperf-sercvice.yml")
	//k8s.ApplyService(*clientset,service)
	////endregion
	//
	////region deployment
	//deployment := k8s.ReadDeploymentYaml("../yamls/hyperf.yaml")
	//_, err := k8s.GetDeploymentByName(*clientset, deployment.Name, deployment.Namespace)
	//if err == nil{
	//	fmt.Printf("deployment %s is exists, it will be deleted",deployment.Name)
	//	k8s.DeleteDeployment(*clientset,deployment)
	//}
	//_, err = k8s.ApplyDeployment(*clientset, deployment)
	//if err != nil {
	//	panic(err)
	//}
	////endregion
	//
}
