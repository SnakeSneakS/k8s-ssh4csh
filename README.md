# (WIP) k8s-ssh4csh 
- SSH (secure shell) 4 (for) CSH (container shell)
- make it possible to access pod through SSH connection
- sidecar pattern (<- I'm not sure this is the best pattern for this usage.)

# required permissions
- check [./manifest/dev/sa.yaml](./manifest/dev/sa.yaml)

# How it works?
1. `kubectl exec` over `ssh` 



# How to develop on k8s?
- check [./manifest/dev/README.md](./manifest/dev/README.md)

# TODO list (note)
- reconsidering the architecture
    - Is sidecar-pattern appropriate for thisr usage?? 
- build & push docker image
- develop injector (using admission webhook)


<!--
# refered web pages while development: 
- [k8ssshの参考になるリポジトリ](https://github.com/guilhem/k8ssh) 
- [k8sのtestの参考 k8s-controller-runtime](https://github.com/kubernetes-sigs/controller-runtime/tree/main)
    - envTest
- [k8sのtestの参考 ginkgo/gomega](https://zenn.dev/zoetro/books/testing-kubernetes-operator/viewer/overview)
- [kubebuilder controller_test](https://zoetrope.github.io/kubebuilder-training/controller-runtime/controller_test.html)
- [In Cluster Config](https://github.com/kubernetes/client-go/tree/master/examples/in-cluster-client-configuration)
- [github jun06t/kubernetes-sample/envoy-service-mesh](https://github.com/jun06t/kubernetes-sample/tree/1f935c66441e9a6fc2211ac66cf11d3af3d341cd/envoy-service-mesh)
- [microsoft sidecar pattern](https://learn.microsoft.com/ja-jp/azure/architecture/patterns/sidecar)
-->
