# kubegenie

### kubeadm、kubelet、kubectl 下载地址
https://storage.googleapis.com/kubernetes-release/release/${RELEASE}/bin/linux/amd64/{kubeadm,kubelet,kubectl}

https://kubernetes-release.pek3b.qingstor.com/release/${RELEASE}/bin/linux/${ARCH}/kubeadm

.
+-- kubernetes  
|   +-- ${version}  
|   |   +-- ${arch}  
|   |   |   +-- kubeadm  
|   |   |   +-- kubelet  
|   |   |   +-- kubectl  
+-- libs
|   +-- rpms  
|   |   +-- containerd.io.rpm  
|   |   +-- docker-ce.rpm  
|   |   +-- docker-ce-cli.rpm  
+-- images  
|   +-- kubernetes  
|   |   +-- ${version}  
|   |   |   +-- k8s.gcr.io/kube-apiserver:v1.19.0  
|   |   |   +-- k8s.gcr.io/kube-controller-manager:v1.19.0   
|   |   |   +-- k8s.gcr.io/kube-scheduler:v1.19.0  
|   |   |   +-- k8s.gcr.io/kube-proxy:v1.19.0  
|   |   |   +-- k8s.gcr.io/pause:3.2  
|   |   |   +-- k8s.gcr.io/etcd:3.4.9-1  
|   |   |   +-- k8s.gcr.io/coredns:1.7.0  
|   |   |   +-- push-images.sh    

#### 配置示例
```
kind: ""
apiversion: ""
masters:
  - 10.0.2.5
workers: []
sshAuth:
    username: root
    password: "123456"
    privateKeyPath: ""
    privateKeyPassword: ""
kubernetes:
    version: v1.20.0
    imageRepo: ""
    apiServerAddress: 10.0.2.5
    apiServerCertSANs: []
    nodeCidrMaskSize: 24
    maxPods: 110
network:
    podCIDR: 172.16.0.1/16
    serviceCIDR: 192.168.0.1/16
    dnsDomain: cluster.local
    calico:
        version: v3.8.2
        ipipMode: Always
        vethMTU: 1440
vip: 10.0.0.255
registries:
    registryMirrors: []
    insecureRegistries: []
    privateRegistry: ""
pkgPath: ""
```

#### 启动命令
```
kubegenie start --config [configpath]
```

kubegenie 将安装步骤解耦，每一步都可以单独执行，不熟悉安装步骤的初学者可以一步一步执行去理解安装过程：
```
# kubegenie start phase -h
Use this command to invoke single phase of the start workflow

Usage:
  kubegenie start phase [command]

Available Commands:       
  initos          
  initpackage    
  installdocker  
  initkubelet      
  initkubeadmconfig           
  initmaster0        
  joinmasters       
  initmasters       
  joinworkers    
  initcluster  
  initcalico       

Flags:
  -h, --help   help for phase

Use "kubegenie start phase [command] --help" for more information about a command.
```

#### kubernetes v1.20.0 离线包
离线包链接：https://pan.baidu.com/s/17KwtKV_AgYuqq5frk98Ezg 提取码: 3tbt 
