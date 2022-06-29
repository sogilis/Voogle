# ADR-013-transformation-services-infra

* Creation Date: 28/06/22
* Status: Pending

## Context

As we start our deployment on Squarescale, we have to start thinking about the scalability and the dependencies of our services. For now, our API cannot work without our transformation services ready. We wish to separate them so they can launch without one another and recover if their counterpart crash. Also, we can add new transformation services without restarting the entire app.

Adress are handled by Squarescale so we have to rely on Consul to fetch the list of existing services or have the services announce themselve at launch.

## Decision


## Options - Services registration

### 1. Consul take the lead

Consul is in charge of networking automation. As such, we may be able to add an event listener and get it to send an updated list of the transformation services when a new transformation service is started/stopped.

#### Benefits
- Act only when needed.
- Detection of a service absence/crash is left to Consul.
- We don't need to reconnect to the service with each request.

#### Drawbacks
- Not sure it's Consul task.
- Need a cache to keep the list in the API.
- In case of multiple API, Consul may update only one per request.
- Harder implementation.

### 2. API ask Consul

The API is managing the user requests. We could fetch the list of existing transformation/services for each transformation request, because Squarescale register new service automatically on Consul. 

#### Benefits
- Ensure informations are up to date when needed.
- Detection of a service absence/crash is left to Consul.
- We don't need to cache the adress/client.
- Simple implementation.

#### Drawbacks
- APIs will make lots of call.
- We will need to recreate a RPC client for each request.

#### Points of attention
We can also maintain a cache to keep the list in the API. Then, we can reduce API calls to Consul and service reconnection (may increase time to detect fault). 

### 3. Services ask API

Another way to update the list of services is by letting them annonce themselves on an API endpoint when they are up and running. To do so, the API need maintain a list of running services.

#### Benefits
- Remove calls between API and Consul.
- We don't need to recreate a RPC client for each request.
- Simple implementation.

#### Drawbacks
- No warning when a service crash.
- In case of multiple API, the service will only contact one.
- Need to keep the list of running services in the API.
