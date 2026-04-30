# Forward packets through the host

## Prepare NFT with masquerade 
```bash
#!/usr/local/sbin/nft -f

#
# flush ruleset
#
# Clear the whole ruleset.
#

flush ruleset

#
# add tables
#

# create empty filter table for IPv4.
# This table is used for standard packet filtering (accept, drop, reject, log, etc.) on forward, input,
# and output chains.
add table ip filter

# create empty nat table for IPv4.
# This table is used for network address translation (masquerade, source NAT, destination NAT) and typically
# contains chains like prerouting, postrouting, and output.
add table ip nat

#
# add chains
#

# create base chain named forward in filter table for IPv4. It hooks into network stack's forward path
# (traffic passing through the host) with default priority 0.
add chain ip filter forward { type filter hook forward priority 0; }

# create base chain named input in filter table for IPv4. It hooks into input path (traffic destined
# for local processes) with priority 0.
add chain ip filter input { type filter hook input priority 0; }

# create base chain named output in filter table for IPv4. It hooks into output path (locally generated
# traffic leaving the host) with priority 0.
add chain ip filter output { type filter hook output priority 0; }

# create base chain named postrouting in nat table for IPv4. It hooks into postrouting stage (after routing,
# just before packets leave an interface) with priority 100 (typical for NAT operations).
add chain nat postrouting { type nat hook postrouting priority 100 ; }

#
# add rules to the forward chain
#

# accept and count packets arriving on interface ens3 and leaving on tun0.
# Used to allow forwarded traffic from the physical interface to the tunnel.
add rule ip filter forward iifname ens3 oifname tun0 counter accept

# accept and count packets arriving on tun0 and leaving on ens3.
# Complements the previous rule, allowing return traffic from the tunnel (tun0) to the physical interface (ens3).
add rule ip filter forward iifname tun0 oifname ens3 counter accept

# add rule to postrouting chain in the nat table: dynamically changes the source IP of packets coming from
# tun0 and heading out ens3 to the IP address of ens3. This hides the tunnel’s internal addresses behind the host’s
# external interface (classic source NAT).
add rule ip nat postrouting iifname tun0 oifname ens3 masquerade
```
-----
## Enable packet forwarding

Run as **root** or with **sudo**

```bash
sysctl -w net.ipv4.ip_forward=1 # enable packet forwarding, disabled by default
```

To persist the command above open `/etc/sysctl.conf`, find and uncommend next line `net.ipv4.ip_forward=1`
