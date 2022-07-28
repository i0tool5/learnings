<h1 style="text-align: center;"> tar 🤝 ssh </h1>

---

After executing command below, the file `file.tar.gz` will appear in `/destination` directory on client host with contents of `/dir/` from the remote host

```sh
ssh user@box tar czf - /dir1/ > /destination/file.tar.gz
```
---

```sh
tar zcvf - /wwwdata | ssh user@backups_server "cat > /backup/wwwdata.tar.gz"
```
The command above will create the `wwwdata.tar.gz` file in the `/backup` directory on the remote host (**backups_server**)

---

```sh
ssh user@server 'tar zcf - /some/dir' | tar zxf -
```
This command is useful when you need to copy the remote directory, which isn't archived yet

---

The commands with sudo can produce *`sudo: sorry, you must have a tty to run sudo`* error. To avoid this, you need to pass the **`-t`** flag to the ssh command

```sh
tar zcvf - /wwwdata | ssh -t user@backups_server "sudo -t cat > /backup/wwwdata.tar.gz"
# NOTE: sudo with the -t flag
```

> man 1 ssh: -t flag - force pseudo-terminal allocation. This can be used to execute arbitrary screen-based programs on a remote machine, which can be very useful, e.g. when implementing menu services.  Multiple -t options force tty allocation, even if ssh has no local tty.

---
The following command can be used to create backup of the entire drive:
```sh
dd if=/dev/sdvf | ssh user@backups_server 'dd of=prod-disk-hostname-sdvf-01-01-1970.img' 
```
To restore a local drive from the image on the server, reverse the command
```sh
ssh user@backups_server 'dd if=prod-disk-hostname-sdvf-01-01-1970.img' | dd of=/dev/sdvf
```

---

Sometimes it is necessary to copy data from the one system to the another. The problem with `scp` and other commands copying the directory structure is that *symbolic links*, *special devices*, *sockets*, *named pipes*, and other stuff **not copied**. Hence, we use tar over ssh.

```ssh
ssh user@old-server 'tar czf - /home/user' | tar xvzf - -C /home/user 
```

---

To tar over SSH with progress bar `pv` tool can be used
```sh
tar zcf - . | pv | ssh user@backups_server "cat > /backups/backup.tgz"
```
> Note: The pv command may not be installed. Use your distribution's package manager to install it.

### Conclusion

The tar and ssh commands are very powerful by themselves. But combine them together gives you superpower 💪 and allows to make some tasks easier

## Links

- [Original article](https://www.cyberciti.biz/faq/howto-use-tar-command-through-network-over-ssh-session/)
- Man7.org [SSH](https://man7.org/linux/man-pages/man1/ssh.1.html)
- Man7.org [tar](https://man7.org/linux/man-pages/man1/tar.1.html)
- Man7.org [pv](https://man7.org/linux/man-pages/man1/pv.1.html)