# Routing protocols: OSPF and IS-IS

**OSPF** (open shortest path first) and **IS-IS** (Intermediate System - Intermediate System) are interior gateway (IGPs) link-state routing protocols. Link-state routing protocols are based on distributed map of the network. Each router that run a link-state protocol have the same copy of this *"network map"*.

## OSPF

To distribute information about the network, OSPF uses **link state advertisements** (*LSA*). Link state advertisements carries information about routes, which is flooded among the other routers. This process is used to create the complete map of the network and save it into link state database. Information stored in link state database must be completely identical on all of the routers in same OSPF area. Currently two versions of OSPF are used: OSPFv2 for IPv4 and OSPFv3 for IPv6.

### OSPF areas

OSPF (since v2) using areas to limit the size of the link-state database, the amount of flooded information and the time that takes to run shortest path algorithm. OSPF area is a logical grouping of routers sharing the same area id. Area ID is similar to IPv4 address and the maximum value is limited to 32 bits. Position of the router in the network topology with respect to OSPF area is important. 

**`Types of link state advertisements (LSA)`**:
- Router LSA (Type1)
- Network LSA (Type2)
- Network summary LSA (Type3)
- ASBR summary (Type4)
- External summary (Type5)
- Type7 is actually the full analog of LSA Type5 for NSSA Area types. When crossing the Area boundary, they are turned into them.

**`Types of the area`**:
- **Backbone** - special area, which always has ID set to 0.0.0.0 (Area 0). Its specialty is that this area is a central node in the OSPF routing domain. Only this area can generate summary routing topology information used by other areas. The link information of other areas is transmitted through Area 0. This also means that all other areas must be connected to Area 0 over Area Border Router. This area supports Type1, Type2, Type3, Type4, and Type5 LSAs.
- **Non-Backbone Non-Stub** (*NBNS*) - also known as Standard Area, which can be considered a smaller version of Area 0. In this area OSPF packets can be normally transmitted and the only diffirence between backbone and non-backbone non-stub is that this area can't be Area 0. This area supports Type1, Type2, Type3, Type4, and Type5 LSAs.
- **Stub** -  this area does not accept any external routes of non-OSPF (can't have links outside autonomous system). If it wants to reach external routes information, it should obtain it through the area border router. This area supports Type1, Type2 and Type3 LSAs.
- **Total Stub Area** - also called completely stub area or totally stubby area.This area does not accept external routes and does not accept the link information of other areas outside of their own area. If it wants to reach the target network outside the area, it will send out the message through the area border router just like the stub area. It should be noted here that since the default route is sent using Type 3 LSAs, the completely stub area does not allow ordinary Type 3 LSA packets, but it supports this type of LSA with default routes. This area supports Type 1, Type 2 LSAs, and Type 3 LSAs with default routes.
- **Not So Stubby Area** (*NSSA*) - this area actually derived from the stub area, it means that in the case of a stub area, it also has ability to send external routes to the other areas over autonomous system border router.

## IS-IS

Long story short: IS-IS is an OSPF-like routing protocol implementation made by ISO. This protocol is not suited for TCP/IP by default, but it has been adapted for IP in order to carry IP routing information inside non-IP packets.

## SIMILARITIES OF OSPF AND IS-IS
- Both IS-IS and OSPF are link-state protocols that maintain a link-state database and run an SPF algorithm based on Dijkstra to compute a shortest path tree of routes.
- Both use Hello packets to create and maintain adjacencies between neighboring routers.
- Both use areas that can be arranged into a two-level hierarchy or into interarea and intraarea routes.
- Both can summarize addresses advertised between their areas.
- Both are classless protocols and handle VLSM.
- Both will elect a designated router on broadcast networks, although IS-IS calls it a designated intermediate system (DIS).
- Both can be configured with authentication mechanisms.

# Links
- Router Switch about [OSPF area types](https://www.router-switch.com/faq/five-ospf-area-types.html)
- Good explanation of [OSPF LSA types](https://networklessons.com/ospf/ospf-lsa-types-explained)
