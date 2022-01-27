# ADR-007-video-processing

* Creation Date: 27/01/2022
* Status: Accepted

## Context
We are able to upload from our Webapp to our S3 bucket. But only half of the job is done. Now we need to transform those videos to the HLS format.

## Decision

We will use the option 1 with a Redis as an event bus. To have more details, see [this document.](./../archi/architecture.md)

## Options
### 1. Processing service with an event bus
Processing videos is an expensive task for the CPU and the GPU. So parallelling videos processing on the same instance can have a negative effect.

We could have several small processing unit that only process one video at the time.


And we can use Protobuf or a similar technology to make our communication more reliable, explicit, and it's supportable by almost all languages.
#### Benefits
* Distributed
* Asynchronous
* Resilient and cross-language
#### Drawbacks
* Scaling of the processing services can be tedious

### 2. Direct communication between services
Services communicated synchronously. It's easier to implement but it's less resilient, especially if the processing servers are busy
#### Benefits
* Easier implementation
#### Drawbacks
* Not fault-tolerant

## Technical resources
- [Redis as an Event Store](https://redis.com/blog/use-redis-event-store-communication-microservices/)
