# ADR-013-services-discovery-cache

* Creation Date: 20/07/22
* Status: Accepted

## Context
We previously chose to use the Consul SDK to retrieve transformation services addresses. For now, each time we need to find a transformation service address, we ask consul.

## Problem
We are doing a lot of request to Consul. If 1 client ask for a video transformed by sevice A and B, we will observe 2 requests to Consul for each video part (api search for A, A search for B).
Considering that a hls video part is commonly between 5 and 10 seconds, we will observe between 12 and 24 requests for a one minute video.
Now multiply that by the number of clients and you obtain possibly thousands of useless request to Consul.
So, we need to save the services addresses. Mostly, we need to be up to date with Consul, and then, we need to find the perfect momement/event to update the saved services addresses cache.

## Decision
We chose to use Consul Watches because it could be very interesting. Moreover, squarescale register/deregister services on startup/crash.

## Options
### 1. Update cache periodically
We can periodically ask Consul to get the list of running services, and store the list locally. Note that each transformation service and api instances will maintains their own cache. To do so, our service discovery can launch a worker in charge of request Consul each 10 seconds for example.

#### Benefits
  - Rather simple
  - Will detect new service AND service crash
  - We can use local API cache to display available services on front

#### Drawbacks
  - No real consistency
  - Still useless request to Consul

### 2. Update only when needed
We can fetch addresses on start up and then update the cache only when a serice is not reachable. Means that when there is no addresse for a given request, we update the list by asking Consul. Also, we can remove a service from the cache if we cannot start connection with.

#### Benefits
  - Rather simple
  - Will detect service crash

#### Drawbacks
  - If only one instance of a given service is available in cache, we will detect the instances that start later
  - We cannot use local API cache to display available services on front, we will fetch only known services
  - Always needs to fails one time to update the list if a service is no more reachable

### 3. Update with Consul Watches
Consul offer the the opportunity to be notify when a service is register or deregister. Using the SDK, we are able to subscribe to these events and launch a handler function.

#### Benefits
  - Will detect service crash
  - Will detect service registration
  - Update only one events. We do not need to implement pooling to fetch service list
  - Really interesting feature of Consul to discover
  - We can use local API cache to display available services on front

#### Drawbacks
  - Harder implementation

## Technical resources
https://www.consul.io/docs/dynamic-app-config/watches