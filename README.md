# terraform-kvm-kubespray
Set up HA Kubernetes cluster using KVM, Terraform and Kubespray.

## Getting Started

### Requirements
+ [Git](https://git-scm.com/) 
+ [Cloud-init](https://cloudinit.readthedocs.io/)
+ [Ansible](https://www.ansible.com/) >= v2.6
+ [Terraform](https://www.terraform.io/) **>= v0.13.x**
+ [KVM - Kernel Virtual Machine](https://www.linux-kvm.org/)
+ [Libvirt provider](https://github.com/dmacvicar/terraform-provider-libvirt) - Setup guide is provided in [docs](./docs/libvirt-provider-setup.md).
+ Internet connection on machine that will run VMs and on VMs

*Note: for Terraform v0.12.x see [this branch](https://github.com/MusicDin/terraform-kvm-kubespray/tree/terraform-0.12).*

### Cluster setup

If you haven't yet, [install libvirt provider](docs/libvirt-provider-setup.md).

Move to main directory:
```
cd terraform-kvm-kubespray
```

Change variables to fit your needs:
```
nano terraform.tfvars
```

Variables are set to work out of the box. Only required variables that are not set are: 
+ `vm_image_source` URL or path on file system to OS image,
+ `vm_distro` a Linux distribution of OS image.

**IMPORTANT:** Review variables before initializing a cluster, as current configuration will create 8 VMs which are quite resource heavy!

Execute terraform script:
```bash
# Initializes terraform project
terraform init

# Shows what is about to be done
terraform plan

# Runs/creates project
terraform apply
```

*Note: Installation process can take up to 20 minutes based on a current configuration.*

### Test cluster

All configuration files will be generated in `config/` directory, and one of them will be `admin.conf` which actually is a `kubeconfig` file.
 
Test your cluster by displaying all cluster's nodes:
```
kubectl --kubeconfig=config/admin.conf get nodes
```

### Cluster management

#### Add worker to the cluster

In [terraform.tfvars](./terraform.tfvars) file add *MAC* and *IP* address for a new VM to `vm_worker_macs_ips`. 
  
Execute terraform script to add worker:
```
terraform apply -var 'action=add_worker'
```

#### Remove worker from the cluster

In [terraform.tfvars](./terraform.tfvars) file remove *MAC* and *IP* address of VM that is going to be deleted from `vm_worker_macs_ips`.

Execute terraform script to remove worker:
```
terraform apply -var 'action=remove_worker'
```
#### Upgrade cluster

In [terraform.tfvars](./terraform.tfvars) file modify:
  + `k8s_kubespray_version` and
  + `k8s_version`.
  
Execute terraform script to upgrade cluster:
```
terraform apply -var 'action=upgrade'
```

#### Destroy cluster

To destroy the cluster, simply run:
```
terraform destroy
```

## More documentation
+ [Setup libvirt provider](docs/libvirt-provider-setup.md)
+ [Load balancing](docs/load-balancer.md)
+ Examples: 
    - [Load balancing to ingress controller](docs/examples/lb-and-ingress-controller.md)


## Related projects

If you are looking forward to installing kubernetes cluster on *vSphere* instead of *KVM* check [this project](https://github.com/sguyennet/terraform-vsphere-kubespray).

## Having issues?

In case you have found a bug, or some unexpected behaviour please open an issue.

If you need anything else, you can contact me on GitHub.

## License

[Apache License 2.0](./LICENSE)