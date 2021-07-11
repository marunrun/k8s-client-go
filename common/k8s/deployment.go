package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	"k8s-client-go/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

//type Deployment interface {
//	ReadDeploymentYaml ReadDeploymentYaml
//}

// ReadDeploymentYaml 读取deployment yaml文件
func ReadDeploymentYaml(filename string) appsv1.Deployment {
	return GetDeploymentByJson(common.ParseYamlToJson(filename))
}

// GetDeploymentByJson 通过json转换成deployment
func GetDeploymentByJson(jsonStr []byte) appsv1.Deployment {
	var deployment appsv1.Deployment
	err := json.Unmarshal(jsonStr, &deployment)
	if err != nil {
		panic(err.Error())
	}
	return deployment
}

func ApplyDeployment(client kubernetes.Clientset, deployment appsv1.Deployment) (*appsv1.Deployment, error) {
	var namespace string
	if deployment.Namespace != "" {
		namespace = deployment.Namespace
	} else {
		namespace = "default"
	}

	deploymentsClient := client.AppsV1().Deployments(namespace)
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), &deployment, metav1.CreateOptions{})
	if err != nil {
		return result, err
	}

	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())
	return result, nil
}

func DeleteDeployment(client kubernetes.Clientset, deployment appsv1.Deployment) {
	var namespace string
	if deployment.Namespace != "" {
		namespace = deployment.Namespace
	} else {
		namespace = "default"
	}
	err := client.AppsV1().Deployments(namespace).Delete(context.TODO(), deployment.Name, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("delete deployment %s success!\n", deployment.Name)
}

func GetDeploymentStatus(clientset kubernetes.Clientset, deployment appsv1.Deployment) (success bool, reasons []string, err error) {

	// 获取deployment状态
	k8sDeployment, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), deployment.Name, metav1.GetOptions{})
	if err != nil {
		return false, []string{"get deployments status error"}, err
	}

	// 获取pod的状态
	labelSelector := ""
	for key, value := range deployment.Spec.Selector.MatchLabels {
		labelSelector = labelSelector + key + "=" + value + ","
	}
	labelSelector = strings.TrimRight(labelSelector, ",")

	podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
	if err != nil {
		return false, []string{"get pods status error"}, err
	}

	readyPod := 0
	unavailablePod := 0
	var waitingReasons []string
	for _, pod := range podList.Items {
		// 记录等待原因
		for _, containerStatus := range pod.Status.ContainerStatuses {
			if containerStatus.State.Waiting != nil {
				reason := "pod " + pod.Name + ", container " + containerStatus.Name + ", waiting reason: " + containerStatus.State.Waiting.Reason
				waitingReasons = append(waitingReasons, reason)
			}
		}

		podScheduledCondition := GetPodCondition(pod.Status, corev1.PodScheduled)
		initializedCondition := GetPodCondition(pod.Status, corev1.PodInitialized)
		readyCondition := GetPodCondition(pod.Status, corev1.PodReady)
		containersReadyCondition := GetPodCondition(pod.Status, corev1.ContainersReady)

		if pod.Status.Phase == "Running" &&
			podScheduledCondition.Status == "True" &&
			initializedCondition.Status == "True" &&
			readyCondition.Status == "True" &&
			containersReadyCondition.Status == "True" {
			readyPod++
		} else {
			unavailablePod++
		}
	}

	// 根据container状态判定
	if len(waitingReasons) != 0 {
		return false, waitingReasons, nil
	}

	// 根据pod状态判定
	if int32(readyPod) < *(k8sDeployment.Spec.Replicas) ||
		int32(unavailablePod) != 0 {
		return false, []string{"pods not ready!"}, nil
	}

	// deployment进行状态判定
	availableCondition := GetDeploymentCondition(k8sDeployment.Status, appsv1.DeploymentAvailable)
	progressingCondition := GetDeploymentCondition(k8sDeployment.Status, appsv1.DeploymentProgressing)

	if k8sDeployment.Status.UpdatedReplicas != *(k8sDeployment.Spec.Replicas) ||
		k8sDeployment.Status.Replicas != *(k8sDeployment.Spec.Replicas) ||
		k8sDeployment.Status.AvailableReplicas != *(k8sDeployment.Spec.Replicas) ||
		availableCondition.Status != "True" ||
		progressingCondition.Status != "True" {
		return false, []string{"deployments not ready!"}, nil
	}

	if k8sDeployment.Status.ObservedGeneration < k8sDeployment.Generation {
		return false, []string{"observed generation less than generation!"}, nil
	}

	// 发布成功
	return true, []string{}, nil
}

func GetDeploymentByName(client kubernetes.Clientset, deploymentName string, namespace string) (*appsv1.Deployment, error) {
	if namespace == "" {
		namespace = "default"
	}

	deployment, err := client.AppsV1().Deployments(namespace).Get(context.TODO(), deploymentName, metav1.GetOptions{})
	if err != nil {
		return deployment,err
	}

	return deployment,nil
}

func PrintDeploymentStatus(clientset kubernetes.Clientset, deployment appsv1.Deployment) {

	// 拼接selector
	labelSelector := ""
	for key, value := range deployment.Spec.Selector.MatchLabels {
		labelSelector = labelSelector + key + "=" + value + ","
	}
	labelSelector = strings.TrimRight(labelSelector, ",")

	for {
		// 获取k8s中deployment的状态
		k8sDeployment, err := clientset.AppsV1().Deployments("default").Get(context.TODO(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Println(err)
		}

		// 打印deployment状态
		fmt.Printf("-------------deployment status------------\n")
		fmt.Printf("deployment.name: %s\n", k8sDeployment.Name)
		fmt.Printf("deployment.generation: %d\n", k8sDeployment.Generation)
		fmt.Printf("deployment.status.observedGeneration: %d\n", k8sDeployment.Status.ObservedGeneration)
		fmt.Printf("deployment.spec.replicas: %d\n", *(k8sDeployment.Spec.Replicas))
		fmt.Printf("deployment.status.replicas: %d\n", k8sDeployment.Status.Replicas)
		fmt.Printf("deployment.status.updatedReplicas: %d\n", k8sDeployment.Status.UpdatedReplicas)
		fmt.Printf("deployment.status.readyReplicas: %d\n", k8sDeployment.Status.ReadyReplicas)
		fmt.Printf("deployment.status.unavailableReplicas: %d\n", k8sDeployment.Status.UnavailableReplicas)
		for _, condition := range k8sDeployment.Status.Conditions {
			fmt.Printf("condition.type: %s, condition.status: %s, condition.reason: %s\n", condition.Type, condition.Status, condition.Reason)
		}

		// 获取pod状态
		podList, err := clientset.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{LabelSelector: labelSelector})
		if err != nil {
			fmt.Println(err)
			return
		}

		for index, pod := range podList.Items {
			// 打印pod的状态
			fmt.Printf("-------------pod %d status------------\n", index)
			fmt.Printf("pod.name: %s\n", pod.Name)
			fmt.Printf("pod.status.phase: %s\n", pod.Status.Phase)
			for _, condition := range pod.Status.Conditions {
				fmt.Printf("condition.type: %s, condition.status: %s, conditon.reason: %s\n", condition.Type, condition.Status, condition.Reason)
			}

			for _, containerStatus := range pod.Status.ContainerStatuses {
				if containerStatus.State.Waiting != nil {
					fmt.Printf("containerStatus.state.waiting.reason: %s\n", containerStatus.State.Waiting.Reason)
				}
				if containerStatus.State.Running != nil {
					fmt.Printf("containerStatus.state.running.startedAt: %s\n", containerStatus.State.Running.StartedAt)
				}
			}
		}

		availableCondition := GetDeploymentCondition(k8sDeployment.Status, appsv1.DeploymentAvailable)
		progressingCondition := GetDeploymentCondition(k8sDeployment.Status, appsv1.DeploymentProgressing)
		if k8sDeployment.Status.UpdatedReplicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.Replicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.AvailableReplicas == *(k8sDeployment.Spec.Replicas) &&
			k8sDeployment.Status.ObservedGeneration >= k8sDeployment.Generation &&
			availableCondition.Status == "True" &&
			progressingCondition.Status == "True" {
			fmt.Printf("-------------deploy status------------\n")
			fmt.Println("success!")
		} else {
			fmt.Printf("-------------deploy status------------\n")
			fmt.Println("waiting...")
		}

		time.Sleep(3 * time.Second)
	}
}

// GetDeploymentCondition returns the condition with the provided type.
func GetDeploymentCondition(status appsv1.DeploymentStatus, condType appsv1.DeploymentConditionType) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func GetPodCondition(status corev1.PodStatus, condType corev1.PodConditionType) *corev1.PodCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}
