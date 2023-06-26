# Samba server

### 1. Install

In some cases **Samba** may be already installed

**Arch** based (Manjaro)
```
$ sudo pacman -S samba
```

**Debian** based (Ubuntu)
```
$ sudo apt install samba
```

### 2. Configure
Create (or choose folder for share)

```
$ mkdir ~/shared
```

Open `/etc/samba/smb.conf` and add the following
```
[Public]
path = /home/<user>/shared
browsable = yes
writable = yes
read only = no
force create mode = 0666
force directory mode = 0777
```
Where \<user\> is a user name home 

### 3. Run Samba

```
$ sudo systemctl restart smb
```

Add user to Samba. User should be created on the host
```
$ sudo smbpasswd -a <username>
```
