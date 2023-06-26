# NFS server

### 1. Install

**Arch** based (Manjaro)
```
$ sudo pacman -S nfs-utils
```

**Debian** based (Ubuntu)
```
$ sudo apt install nfs-kernel-server
```

### 2. Configure
Create (or choose folder for share)
```
$ mkdir ~/shared
```
Change configuration file
```
$ vim /etc/exports
# add the following line, where <user> is a user folder
/home/<user>/shared *(rw,sync,no_subtree_check,no_root_squash)
```
- **\*** - specifies the IP. In this case **\*** means that any client is allowed to access the shared folder
- **rw** - allows both read and write requests. The default is to disallow any request that changes the filesystem. This also can be set explicitly by **ro** parameter
- **sync** - forces NFS to write changes to disk before replying. This results in a more stable and consistent environment since the reply reflects the actual state of the remote volume. However, it also reduces the speed of file operations.
- **no_subtree_check** - disable subtree checking. Subtree cheking is a process where the host must check whether the file is actually still available in the exported tree for every request. This can cause many problems when a file is renamed while the client has it opened. In almost all cases, it is better to disable subtree checking. But this *may have* mild *security implications*.
- **no_root_squash** - mainly useful for diskless clients. By default, NFS translates requests from a root user remotely into a non-privileged user on the server. This was intended as security feature to prevent a root account on the client from using the file system of the host as root. **no_root_squash** disables this behavior for certain shares.

> To learn more about *exports* file, read `man exports`

### 3. Enable and start NFS

**Arch** based (Manjaro)
```
$ sudo systemctl enable nfs-server
$ sudo systemctl start nfs-server
```

**Debian** based (Ubuntu)
```
$ sudo systemctl enable nfs-kernel-server
$ sudo systemctl restart nfs-kernel-server
```

# NFS client

### 1. Install

**Arch** based (Manjaro)
```
$ sudo pacman -S nfs-utils
```

**Debian** based (Ubuntu)
```
$ sudo apt install nfs-common
```

### 2. Configure

Create client mount point
```
mkdir ~/remote_share
```

Mount remote NFS
```
sudo mount -t nfs server_ip:/shared ~/remote_share
```

> Note, on Arch based distros, you need to start (and optionaly enable) nfs client:

```
# On Arch
$ sudo systemctl enable nfs-client.target
$ sudo systemctl restart nfs-client.target
```

## **Optional**: persistent mount

Fill the */etc/fstab* file with the following
```
<server_ip>:/home/<user>/shared ~/remote_share  nfs  defaults,timeo=900,retrans=5   0 0
```

Replace \<server_ip\> and \<user\> with appropriate values.
