# Setting Mikrotik antenna gain

Sometimes you may encounter signal strength problems on mikrotik devices. First you need to know the signal strength of the antenna, which is set by default depending on the region. To do this, open a terminal (WinBox) and enter the following command
```sh
/interfaces/wireless
```
Then type next command, which will show detailed information about wireless interfaces
```sh
print detail advanced
```
In the output search for **antenna gain** parameter. As official documentation says ```Antenna gain in dBi, used to calculate maximum transmit power according to country regulations.```
The higher this parameter is, then weaker the signal strenght, and the opposite.

To change the value of the antena-gain parameter, type into the console:
```sh
set <iface> antenna-gain=<ag>
```
Where *iface* - is an interface name and *ag* desired antenna gain. For example ```set wlan1 antenna-gain=0``` will completely disable transmission power limitation.

# Links

- Official Mikrotik [documentation](https://help.mikrotik.com/docs/display/ROS/Wireless+Interface) for wireless interfaces
