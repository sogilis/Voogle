# ADR-013-services-discovery

* Creation Date: 28/06/22
* Status: Accepted

## Context
Voogle is a micro-services application by design.
As we start our deployment on Squarescale, we have to start thinking about the scalability and the dependencies of our services.

## Problem
To resolve runtime dependencies between services, we need a service discovery that resolve following problems :

 - When service A needs to request service B, the service discovery provides the IP address of service B to service A.
 - When service B starts, it registers to the service discovery. Once registered, service B is available for requests by other services.
 - When service B fails, it is restarded by SquareScale. Service A detects the failure, a new instance of service B is spawn by SquareScale and registered to service discovery, then service A request the service discovery to get the new address of service B.
 - In multi node cluster environment, the service discovery is in charge of gathering all IP addresses of service B. **The service discovery is not in charge of load balancing**.

For now, our API service depends on our transformation services (gray, flip, etc...). Services lifecycle has to be independant and uncoupled, by instance services can start, stop, fail and be restarted independantly.

### Note
Adresses are handled by Squarescale so we have to rely on Consul to fetch the list of existing services or have the services announce themselve at launch. Consul already implements the features we need for our service discovery.

## Decision
We chose to use Consul because Squarescale register automatically our services on it. Also, we will use the Consul SDK which is easier to use than the API and is really made for our use.

## Options
### 1. Retrieving addresses with local DNS request
```go
    address := "gray-server-transformer.service.consul" // Can be retrieve by env var
    conn, err := grpc.Dial(address, opts)
    if err != nil {
        log.Error("Cannot open TCP connection with transformer server :", err)
    }
```

### 2. Directly ask Consul for services addresses, using SDK
```go
   config := &api.Config{
		Address:  consulAddr,
		HttpAuth: BasicAuth,
	}

	// Create a Consul API client
	client, err := api.NewClient(config)
	if err != nil {
		fmt.Println("Cannot create consul client")
	}

	// Create a Consul agent client
	agent := client.Agent()

	services, err := agent.ServicesWithFilter("transformer in Tags and flip in Service")
	if err != nil {
		fmt.Println("Cannot retrieve service")
	}
```

### 3. Directly ask Consul for services addresses, using Rest API
```go
   // Create header
	header := http.Header{}
	header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(user+":"+password)))
	
	// Create url
	parsedURLServices, err := url.Parse("http://" + consulAddr + `/v1/agent/services?filter="transformer"+in+Tags`)
	if err != nil {
		return nil, err
	}

	request: &http.Request{
		Method: http.MethodGet,
		URL:    parsedURLServices,
		Header: header,
	},

	var httpC http.Client
	resp, err := httpC.Do(request)
	if err != nil {
		return TransformerInfos{}, err
	}
```

## Technical resources
https://www.consul.io/api-docs/catalog
https://www.consul.io/api-docs/agent
https://www.consul.io/docs/discovery/services
https://pkg.go.dev/github.com/hashicorp/consul/api