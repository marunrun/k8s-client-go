package main

import (
	"flag"
	"fmt"
	"k8s-client-go/common"
	"k8s-client-go/common/k8s"
)

func main() {

	//resource := flag.String("resource", "", "需要操作的资源,[deployment,pod,service]")
	//action := flag.String("action", "", "动作,[]")

	flag.Parse()

	// 初始化k8s客户端
	clientset := common.GetK8sClient()

	//region service
	service := k8s.ReadServiceYaml("D:\\code\\yamls\\hyperf-sercvice.yml")
	k8s.ApplyService(*clientset,service)
	//endregion

	//region deployment
	deployment := k8s.ReadDeploymentYaml("D:\\code\\yamls\\hyperf.yaml")
	_, err := k8s.GetDeploymentByName(*clientset, deployment.Name, deployment.Namespace)
	if err == nil{
		fmt.Printf("deployment %s is exists, it will be deleted",deployment.Name)
		k8s.DeleteDeployment(*clientset,deployment)
	}
	_, err = k8s.ApplyDeployment(*clientset, deployment)
	if err != nil {
		panic(err)
	}
	//endregion
}

