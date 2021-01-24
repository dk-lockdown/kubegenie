package tmpl

import (
	"text/template"
)

import (
	"github.com/lithammer/dedent"
)

var DockerDaemonJsonTmpl = template.Must(template.New("DockerConfig").Parse(
	dedent.Dedent(`{
  "log-opts": {
    "max-size": "5m",
    "max-file":"3"
  },
  {{- if .Mirrors }}
  "registry-mirrors": [{{ .Mirrors }}],
  {{- end}}
  {{- if .InsecureRegistries }}
  "insecure-registries": [{{ .InsecureRegistries }}],
  {{- end}}
  "exec-opts": ["native.cgroupdriver=systemd"]
}
    `)))

var KubeadmCfgTmpl = template.Must(template.New("kubeadmCfg").Parse(
	dedent.Dedent(`---
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
imageRepository: {{ .ImageRepo }}
kubernetesVersion: {{ .Version }}
certificatesDir: /etc/kubernetes/pki
controlPlaneEndpoint: {{ .APIServerAddress }}:6443
networking:
  dnsDomain: {{ .DNSDomain }}
  podSubnet: {{ .PodCIDR }}
  serviceSubnet: {{ .ServiceCIDR }}
apiServer:
  extraArgs:
    anonymous-auth: "true"
    bind-address: 0.0.0.0
    insecure-port: "0"
    profiling: "false"
    apiserver-count: "1"
    endpoint-reconciler-type: lease
    authorization-mode: Node,RBAC
    enable-aggregator-routing: "false"
    allow-privileged: "true"
    storage-backend: etcd3
    audit-log-maxage: "30"
    audit-log-maxbackup: "10"
    audit-log-maxsize: "100"
    audit-log-path: /var/log/apiserver/audit.log
    feature-gates: CSINodeInfo=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,RotateKubeletClientCertificate=true
  certSANs:
    {{ range .CertSANs -}}
    - {{ . }}
    {{ end -}}
	- {{.VIP}}
controllerManager:
  extraArgs:
    node-cidr-mask-size: "{{ .NodeCidrMaskSize }}"
    experimental-cluster-signing-duration: 87600h
    bind-address: 127.0.0.1
    profiling: "false"
    port: "10252"
    terminated-pod-gc-threshold: "10"
    feature-gates: RotateKubeletServerCertificate=true,CSINodeInfo=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,RotateKubeletClientCertificate=true
  extraVolumes:
  - name: host-time
    hostPath: /etc/localtime
    mountPath: /etc/localtime
    readOnly: true
scheduler:
  extraArgs:
    profiling: "false"
    bind-address: 127.0.0.1
    port: "10251"
    feature-gates: CSINodeInfo=true,VolumeSnapshotDataSource=true,ExpandCSIVolumes=true,RotateKubeletClientCertificate=true

---
apiVersion: kubeproxy.config.k8s.io/v1alpha1
kind: KubeProxyConfiguration
bindAddress: 0.0.0.0
clientConnection:
 acceptContentTypes: 
 burst: 10
 contentType: application/vnd.kubernetes.protobuf
 kubeconfig: 
 qps: 5
clusterCIDR: {{ .PodCIDR }}
configSyncPeriod: 15m0s
conntrack:
 maxPerCore: 32768
 min: 131072
 tcpCloseWaitTimeout: 1h0m0s
 tcpEstablishedTimeout: 24h0m0s
enableProfiling: False
healthzBindAddress: 0.0.0.0:10256
ipvs:
 excludeCIDRs: 
 - "{{.VIP}}/32"
 minSyncPeriod: 0s
 scheduler: rr
 syncPeriod: 30s
 strictARP: False
mode: "ipvs"

---
apiVersion: kubelet.config.k8s.io/v1beta1
kind: KubeletConfiguration
maxPods: {{ .MaxPods }}
rotateCertificates: true
kubeReserved:
  cpu: 200m
  memory: 250Mi
systemReserved:
  cpu: 200m
  memory: 250Mi
evictionHard:
  memory.available: 5%
evictionSoft:
  memory.available: 10%
evictionSoftGracePeriod: 
  memory.available: 2m
evictionMaxPodGracePeriod: 120
evictionPressureTransitionPeriod: 30s
featureGates:
  CSINodeInfo: true
  VolumeSnapshotDataSource: true
  ExpandCSIVolumes: true
  RotateKubeletClientCertificate: true

    `)))

var KubeletServiceTmpl = template.Must(template.New("kubeletService").Parse(
	dedent.Dedent(`[Unit]
Description=kubelet: The Kubernetes Node Agent
Documentation=http://kubernetes.io/docs/

[Service]
ExecStart=/usr/local/bin/kubelet
Restart=always
StartLimitInterval=0
RestartSec=10

[Install]
WantedBy=multi-user.target
    `)))

var KubeletEnvTmpl = template.Must(template.New("kubeletEnv").Parse(
	dedent.Dedent(`# Note: This dropin only works with kubeadm and kubelet v1.11+
[Service]
Environment="KUBELET_KUBECONFIG_ARGS=--bootstrap-kubeconfig=/etc/kubernetes/bootstrap-kubelet.conf --kubeconfig=/etc/kubernetes/kubelet.conf"
Environment="KUBELET_CONFIG_ARGS=--config=/var/lib/kubelet/config.yaml"
# This is a file that "kubeadm init" and "kubeadm join" generate at runtime, populating the KUBELET_KUBEADM_ARGS variable dynamically
EnvironmentFile=-/var/lib/kubelet/kubeadm-flags.env
# This is a file that the user can use for overrides of the kubelet args as a last resort. Preferably, the user should use
# the .NodeRegistration.KubeletExtraArgs object in the configuration files instead. KUBELET_EXTRA_ARGS should be sourced from this file.
EnvironmentFile=-/etc/default/kubelet
Environment="KUBELET_EXTRA_ARGS=--node-ip={{ .NodeIP }}  --hostname-override={{ .Hostname }}"
ExecStart=
ExecStart=/usr/local/bin/kubelet $KUBELET_KUBECONFIG_ARGS $KUBELET_CONFIG_ARGS $KUBELET_KUBEADM_ARGS $KUBELET_EXTRA_ARGS
    `)))
