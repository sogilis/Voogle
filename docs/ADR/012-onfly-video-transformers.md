# ADR-012-onfly-video-transformers

* Creation Date: 15/06/2022
* Status: Accepted

## Context
We want that the user can select 2 or 3 video transformations that will be apply on reading video. We also want one service for each transformation. These services should transform each HLS video parts by the time the user ask for them. The user can freely ask for no transformation, one transformation or every composition of transformations.

To do so, we need to choose a means to communicate between api, transformations services and the video storage service.

The solution must be suitable for the principal purpose of Voogle : read video streams. It means that the process of video transformation must be as fast as possible, and since we transfer video part, we have to limit the network load. 

## Decision
In order to be a Squarescale demonstrator, Voogle cannot use AWS mediaconvert. Moreover we already use and deploy a service with FFMPEG, it's the easiest solution. We will use 
For the services communication, we need a synchronous communication. In order to develop team skills, and because it seems to be easily maintainable especially to add new transformation services in the future, we will use RPC protocol with gRPC.

## Options - Video transformation process

### 1. FFMPEG
For the video transformation, we could use ffmpeg

#### Benefits
  - We already use it for the video encoding
  - We have knwoledge on it
  - "Free" solution
  - Easy deployment

#### Drawbacks
  - Not easy to set up
  - Fast process, but need a powerful machine
  - Greddy process in memory and cpu

### 2. AWS Mediaconvert
AWS offer a solution to encode and process videos

#### Benefits
  - No more memory/cpu issues
  - Very fast process

#### Drawbacks
  - Hard deployment
  - Paid service
  - Not using Squarescale

## Options - Services communication

### 1. REST API
We could use simple HTTP requests with a REST API. (Synchronous communication)

#### Benefits
- Simple architecture
- Same as our API
- Low latency
- Streams friendly

#### Drawbacks
- Strong coupling
- Cumulative latency
- Harder scalability

### 2. RPC
The remote procedure call protocol is another way to deal with Synchronous communication.

#### Benefits
- Low latency
- Streams friendly
- Easily maintainable
- Develop skills

#### Drawbacks
- Harder to set up
- Strong coupling
- Cumulative latency
- Harder scalability

### 3. Event Bus
We could use an event bus, with a publish/subscribe. (Asynchronous communication)

#### Benefits
  - High scalability
  - Decoupling

#### Drawbacks
  - High latency

#### Points of attention
Because the user is actually waiting for a video part, the high latency is prohibitive.

## Technical resources
- [Architecture Decision Record](https://github.com/joelparkerhenderson/architecture-decision-record/blob/main/examples/programming-languages/index.md)
- [AWS Mediaconvert](https://aws.amazon.com/fr/mediaconvert/)
- [FFMPEG](https://ffmpeg.org/)
- [Formation Cloud Native](https://sogilis.com/)
- [Communication microservices](https://docs.microsoft.com/fr-fr/dotnet/architecture/microservices/architect-microservice-container-applications/communication-in-microservice-architecture)
