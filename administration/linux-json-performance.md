This is an excerpt from an [article](https://talawah.io/blog/extreme-http-performance-tuning-one-point-two-million/) on how to improve the performance of HTTP requests handling on Linux. The original article says that these methods made it possible to increase the performance of a JSON processor based on the *libreactor* library in an **Amazon EC2 environment (4 vCPU)** from 224 thousand API requests per second with stock settings of **Amazon Linux 2 with a 4.14 kernel** to 1.2 million requests per second after optimization ( an increase of 436%), and also led to a reduction in delays in processing requests by 79%.

> **WARNING**: these optimizations are not a panacea. Use them at your own risk

## Main optimizations
- ***libreactor* code optimization**. The R18 variant from the Techempower set was used as a basis, which was improved by removing the code to limit the number of CPU cores used (optimization made it possible to speed up work by 25-27%), building in GCC with the "-O3" options (increase 5-10% ) and "-march-native" (5-10%), replacing read/write calls with recv/send (5-10%) and reducing overhead when using pthreads (2-3%). The overall performance gain after code optimization was 55%, and throughput increased from 224k req/s to 347k req/s.

- **disable protection against vulnerabilities caused by speculative instruction execution**. Using kernel boot options "nospectre_v1 nospectre_v2 pti=off mds=off tsx_async_abort=off" allowed performance to be increased by 28% and throughput increased from 347k req/s to 446k req/s. Separately, the increase from the parameter "nospectre_v1" (protection against Specter v1 + SWAPGS) was 1-2%, "nospectre_v2" (protection against Specter v2) - 15-20%, "pti=off" (Spectre v3/Meltdown) - 6 %, "mds=off tsx_async_abort=off" (MDS/Zombieload and TSX Asynchronous Abort) - 6%. Left unchanged settings to protect against attacks L1TF/Foreshadow (l1tf=flush), iTLB multihit, Speculative Store Bypass and SRBDS, which did not affect performance, as they did not intersect with the tested configuration (for example, specific to KVM, nested virtualization and others CPU models).

- **disable auditing and system call blocking mechanisms** with the "auditctl -a never,task" command and specifying the "--security-opt seccomp=unconfined" option when starting the docker container. The overall performance gain was 11% and throughput increased from 446k req/s to 495k req/s.

- **disabling iptables/netfilter by unloading their associated kernel modules**. The idea to disable a firewall that was not used in a specific server solution was prompted by profiling results, judging by which the nf_hook_slow function took 18% of the time to execute. It is noted that nftables is more efficient than iptables. After disabling iptables, there was a 22% performance gain and throughput increased from 495k req/s to 603k req/s.

- **reduce the migration of handlers between different CPU cores to improve the efficiency of using the processor cache**. Optimization was made both at the level of binding libreactor processes to CPU cores (CPU Pinning), and through pinning kernel network handlers (Receive Side Scaling). For example, disabled irqbalance and explicitly set CPU queuing affinities in /proc/irq/$IRQ/smp_affinity_list. To use the same CPU core to process the libreactor process and the network queue of incoming packets, a native BPF handler is used, connected by setting the SO_ATTACH_REUSEPORT_CBPF flag when creating the socket. The /sys/class/net/eth0/queues/tx-<n>/xps_cpus settings have been changed to bind outgoing packet queues to the CPU. The overall performance gain was 38% and throughput increased from 603k req/s to 834k req/s.

- **optimization of interrupt handling and the use of polling**. Enabling adaptive-rx mode in the ENA driver and manipulating the sysctl net.core.busy_read allowed us to increase performance by 28% (bandwidth increased from 834k req / s to 1.06M req / s, and latency decreased from 361μs to 292μs).

- disabling system services that lead to unnecessary locks in the network stack. Disabling dhclient and manually setting the IP resulted in a 6% performance improvement and throughput increased from 1.06M req/s to 1.12M req/s. The reason for the performance impact of dhclient is in analyzing traffic using a raw socket.

- fight with spin lock. Putting the network stack into "noqueue" mode via sysctl "net.core.default_qdisc=noqueue" and "tc qdisc replace dev eth0 root mq" resulted in a performance gain of 2% and throughput increased from 1.12M req/s to 1.15M req/s.

- final minor optimizations such as disabling GRO (Generic Receive Offload) with "ethtool -K eth0 gro off" and replacing the cubic congestion control algorithm with reno with sysctl "net.ipv4.tcp_congestion_control=reno". The overall performance increase was 4%. Throughput increased from 1.15M req/s to 1.2M req/s.

# Links

[Original post]()