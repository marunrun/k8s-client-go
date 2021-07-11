package common

import (
	"istio.io/client-go/pkg/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// GetK8sClient 初始化k8s客户端
func GetK8sClient() *kubernetes.Clientset {

	if environment := os.Getenv("RUNTIME_ENVIRONMENT"); environment == "in_cluster" {
		return getK8sClientInCluster()
	}

	filename := ""
	if home := homeDir(); home != "" {
		filename = filepath.Join(home, ".kube", "config")
	} else {
		panic("can not found kube config")
	}

	return getK8sClientOutCluster(filename)
}

func GetIstioClient() *versioned.Clientset {
	if environment := os.Getenv("RUNTIME_ENVIRONMENT"); environment == "in_cluster" {
		return getIstioClientInCluster()
	}

	filename := ""
	if home := homeDir(); home != "" {
		filename = filepath.Join(home, ".kube", "config")
	} else {
		panic("can not found kube config")
	}

	return getIstioClientOutCluster(filename)
}

// 获取k8s client
func getK8sClientOutCluster(filename string) *kubernetes.Clientset {
	// 生成rest client配置
	restConf, err := clientcmd.BuildConfigFromFlags("", filename)
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(restConf)
	if err != nil {
		panic(err)
	}

	return client
}

// 运行在cluster中
func getK8sClientInCluster() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return client
}
func getIstioClientOutCluster(filename string) *versioned.Clientset {
	restConf, err := clientcmd.BuildConfigFromFlags("", filename)
	if err != nil{
		panic(err)
	}
	client, err := versioned.NewForConfig(restConf)
	return client
}

func getIstioClientInCluster() *versioned.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}
	client, err := versioned.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return client
}


func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
