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
metadata:
    name: ""
    generatename: ""
    namespace: ""
    selflink: ""
    uid: ""
    resourceversion: ""
    generation: 0
    creationtimestamp: "0001-01-01T00:00:00Z"
    deletiontimestamp: null
    deletiongraceperiodseconds: null
    labels: {}
    annotations: {}
    ownerreferences: []
    finalizers: []
    clustername: ""
    managedfields: []
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
    apiservercertsans: []
    nodecidrmasksize: 24
    maxpods: 110
network:
    podcidr: 172.16.0.1/16
    servicecidr: 192.168.0.1/16
    dnsdomain: cluster.local
    calico:
        version: v3.8.2
        ipipmode: Always
        vethmtu: 1440
vip: 10.0.0.255
registries:
    registrymirrors: []
    insecureregistries: []
    privateregistry: ""
pkgPath: ""
```

#### 启动命令
```
kubegenie start --config [configpath]
```

#### kubernetes v1.20.0 离线包

线包链接：https://pan.baidu.com/s/1bdKFSiah2xihFyTGnf9Spg 提取码: db3q