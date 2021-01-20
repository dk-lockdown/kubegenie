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
+-- lib
|   +-- rpm
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
