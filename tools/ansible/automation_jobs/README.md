# automation_jobs

Ansible 自动化任务

## Describe

Ansible playbook，所有项目都应放在一级目录内，调用时通过roles控制

## Getting started

注意，在综合运维管理平台上选择的hosts，最终都会统一归属到`select_hosts`组内，在编写playbook时请注意选用此hosts，如下所示：

```yml
- hosts: select_hosts
  vars:
    gpfs_clients_dns: "{{ hostvars['localhost']['gpfs_clients_dns'].stdout }}"
  gather_facts: false
  remote_user: ansibleuser
  become: yes
  tasks:
    - include_role:
        name: gpfs_modify_hosts
```

所有hosts均需要放在inventory目录的hosts文件内
