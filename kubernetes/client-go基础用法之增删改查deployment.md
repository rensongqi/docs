# 1 client-go简介

`client-go`是一个调用kubernetes集群资源对象API的客户端，即通过client-go实现对kubernetes集群中资源对象（包括deployment、service、ingress、replicaSet、pod、namespace、node等）的增删改查等操作。大部分对kubernetes进行前置API封装的二次开发都通过client-go这个第三方包来实现。

**client-go官方文档**：https://github.com/kubernetes/client-go

# 2 client-go的使用

windows 在家目录下创建一个.kube文件夹，然后把k8s的config文件放入其中，位置：`C:\Users\songqi\.kube\config`

linux 也是需要在go运行的机器上放置k8s访问config文件，位置：`~/.kube/config`

@[TOC]

## 2.1 创建clientSet

```go
package clientset

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func GetClientSet() (*kubernetes.Clientset, error) {
	var kubeConfig *string
	// 从当前系统环境中读取家目录，然后拼接config 路径
	// 或者直接给一个kube config的绝对路径字符串也可
	if home := HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// uses the current context get restConfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		log.Panic(err)
	}

	// 创建clientSet
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return clientSet, nil
}

func HomeDir() string {
	// linux
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	// windows
	return os.Getenv("USERPROFILE")
}
```
## 2.2 获取pod信息

```go
func GetPods(client *kubernetes.Clientset, ctx context.Context, namespace string) {
	// get pod
	pods, err := client.CoreV1().Pods(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("pod Name ===> ", pods.Items[0].Status.ContainerStatuses[0].Name)
	fmt.Println("pod Image ===> ", pods.Items[0].Status.ContainerStatuses[0].Image)
	fmt.Println("pod State ===> ", pods.Items[0].Status.ContainerStatuses[0].State.Running)
}
```
## 2.3 获取deployment信息

```go
func GetDeploy(client *kubernetes.Clientset, ctx context.Context, deployName, namespace string) {
	// get deploy
	deployment, err := client.AppsV1().Deployments(namespace).Get(ctx, deployName, metaV1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	fmt.Println("deployment name ===> ", deployment.Name)
}
```
## 2.4 更新deployment副本数量

```go
func UpdateDeployReplica(client *kubernetes.Clientset, ctx context.Context, deployName, namespace string, replicas int32) {
	// 1 方法一：更新deployment 副本数量
	deployment, err := client.AppsV1().Deployments(namespace).Get(ctx, deployName, metaV1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	// 设置副本数量
	deployment.Spec.Replicas = &replicas
	deployment, err = client.AppsV1().Deployments(namespace).Update(ctx, deployment, metaV1.UpdateOptions{})

	// 2 方法二：更新副本数量的另一种方法
	replica, err := client.AppsV1().Deployments(namespace).GetScale(ctx, deployName, metaV1.GetOptions{})
	replica.Spec.Replicas = replicas
	replica, err = client.AppsV1().Deployments(namespace).UpdateScale(ctx, deployName, replica, metaV1.UpdateOptions{})
	fmt.Println("replica name ====>", replica.Name)
}
```
## 2.5 更新deployment镜像

```go
func UpdateDeployImage(client *kubernetes.Clientset, ctx context.Context, deployName, namespace, image string) {
	deployment, err := client.AppsV1().Deployments(namespace).Get(ctx, deployName, metaV1.GetOptions{})
	if err != nil {
		log.Println(err)
	}
	deployment.Spec.Template.Spec.Containers[0].Image = image
	deployment, err = client.AppsV1().Deployments(namespace).Update(ctx, deployment, metaV1.UpdateOptions{})
}
```
## 2.6 删除deployment

```go
func DeleteDeploy(client *kubernetes.Clientset, ctx context.Context, deployName, namespace string) {
	// 删除deployment
	err := client.AppsV1().Deployments(namespace).Delete(ctx, deployName, metaV1.DeleteOptions{})
	if err != nil{
		log.Println(err)
	}
}
```
## 2.7 创建deployment和service

```go
func CreateDeploy(client *kubernetes.Clientset, ctx context.Context, namespace string) {
	var replicas int32 = 3
	var targetPort int32 = 80
	intString := intstr.IntOrString{
		IntVal: targetPort,
	}
	deployment := &appV1.Deployment{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "nginx",
			Labels: map[string]string{
				"app": "nginx",
			},
			Namespace: namespace,
		},
		Spec: appV1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metaV1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: coreV1.PodTemplateSpec{
				ObjectMeta: metaV1.ObjectMeta{
					Name: "nginx",
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: coreV1.PodSpec{
					Containers: []coreV1.Container{
						{
							Name: "nginx",
							Image: "nginx:1.16.1",
							Ports: []coreV1.ContainerPort{
								{
									Name: "http",
									Protocol: coreV1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	service := &coreV1.Service{
		ObjectMeta: metaV1.ObjectMeta{
			Name: "nginx",
			Labels: map[string]string{
				"app": "nginx",
			},
			Namespace: namespace,
		},
		Spec: coreV1.ServiceSpec{
			Type: coreV1.ServiceTypeNodePort,
			Ports: []coreV1.ServicePort{
				{
					Name: "nginx",
					Port: 80,
					TargetPort: intString,
					NodePort: 30088,
					Protocol: coreV1.ProtocolTCP,
				},
			},
			Selector: map[string]string{
				"app": "nginx",
			},
		},
	}

	deployment, err := client.AppsV1().Deployments(namespace).Create(ctx, deployment, metaV1.CreateOptions{})
	if err != nil {
		log.Println("err ===> ", err)
	}
	service, err = client.CoreV1().Services(namespace).Create(ctx, service, metaV1.CreateOptions{})
}
```

## 2.8 main 函数

```go
package main

import (
	"client-go/clientset"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	appV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

func main() {
	deployName := "myapp-deploy"
	namespace := "default"
	//image := "nginx:1.15-alpine"
	var replicas int32 = 2

	ctx := context.Background()
	// get the clientset
	client, err := clientset.GetClientSet()
	if err != nil {
		log.Panic(err)
	}

	GetPods(client, ctx, namespace)
	GetDeploy(client, ctx, deployName, namespace)
	//DeleteDeploy(client, ctx, deployName, namespace)
	//UpdateDeployImage(client, ctx, deployName, namespace, image)
	UpdateDeployReplica(client, ctx, deployName, namespace, replicas)

	CreateDeploy(client, ctx, namespace)
}
```

**有关client-go的详细说明参考文章：**[client-go的使用及源码分析](https://blog.csdn.net/huwh_/article/details/78821805)