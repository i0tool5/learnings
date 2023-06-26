<h1 style="text-align: center;">Share files and folders over the Network</h1>

- [NFS](./nfs.md) 
- CIFS

## NFS (Network File System)
NFS is a network that was introduced by Sun Microsystems and is used by Unix or Linux-based operating systems. This is a network that is used for giving remote access capabilities to the applications. Remote access enables the user to edit or even take a closer look at his computer by using another computer. Old files can be repaired even when the user is at a distance from his computer. This protocol gives devices the functionality to modify the data over a network. 

## CIFS (Common Internet File System)
CIFS is a Windows-based network for file sharing and is used in devices that run on Windows OS. CIFS was introduced as the public version of Server Message Block (SMB) which was invented by Microsoft. This is a very efficient feature that enables the devices to share multiple devices that are printers and even multiple ports for the user and administration. CIFS also enables a request for accessing files of another computer that is connected to the server. Then this request is served by the server to the requested client. CIFS supports huge data companies to ensure that their data is used by the employees at multiple locations. 

### NFS and CIFS comparsion 

|   |NFS|CIFS|
|---|---|---|
| Communication | Better than CIFS | Creates a confusion when communicating |
| Support | UNIX or LINUX OS are most preferable supported | Windows OS are preferable supported |
| Session feature | Doesn't provide sessions | Provides sessions |
| Port protocols | 111 port for both TCP and UDP | Uses 139 and 455 for TCP and 138 and 137 for UDP |
| Speed & Scalability | Highly scalable and more speed than CIFS | Low scalable and moderate speed |
| Implementation | Simple to implement and executing queries is fast | Difficult to implement and configure for faults |
| Security | Not reliable, no special security | It has more security features than NFS |
| Transportation | NFS is a transport-dependent protocol and offers a high rate of communication speed. | CIFS is generally used for direct hosting and NetBIOS-dependent transport over IP and TCP protocols. |
