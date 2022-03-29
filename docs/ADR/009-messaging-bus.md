# ADR-009-messaging-bus

- Creation Date: 11/02/2022
- Status: Accepted

## Context

The Voogle project is a demonstrator for Squarescale. The software is cloud native project. The services are independent but they must be able to communicate with each other. To carry out this communication, we need an external service of messaging-bus.

## Decision

We decided to use RabbitMQ because this software is really made for our use, no effort, available on AWS, easy to implement. Redis is rather a key storage software, its implementation requires more effort to realize a simple message queue

## Options

### 1. RabbitMQ

RabbitMQ is an open-source message-broker software (sometimes called message-oriented middleware) that originally implemented the Advanced Message Queuing Protocol (AMQP) and has since been extended with a plug-in architecture to support Streaming Text Oriented Messaging Protocol (STOMP), MQ Telemetry Transport (MQTT), and other protocols.

Written in Erlang, the RabbitMQ server is built on the Open Telecom Platform framework for clustering and failover. Client libraries to interface with the broker are available for all major programming languages. The source code is released under the Mozilla Public License.
[Wikipedia](https://en.wikipedia.org/wiki/RabbitMQ)

#### Benefits

- Very easy to implement
- Large community
- Native message queue Go package
- Available on AWS with Amazon MQ service
- Dashboard user
- AMQP ptotocol, AMQP is a protocol, not an implementation
- Persistence management
- Clustering
- Asynchronous sending
- Acknowledgement of receipt

#### Drawbacks

- Not available on Sqsc

### 2. Redis

Redis (Remote Dictionary Server) is an in-memory data structure store, used as a distributed, in-memory key–value database, cache and message broker, with optional durability. Redis supports different kinds of abstract data structures, such as strings, lists, maps, sets, sorted sets, HyperLogLogs, bitmaps, streams, and spatial indices. The project was developed and maintained by Salvatore Sanfilippo. From 2015 until 2020, he led a project core team sponsored by Redis Labs. Salvatore Sanfilippo left Redis as the maintainer in 2020. It is open-source software released under a BSD 3-clause license. In 2021, not long after the original author and main maintainer left, Redis Labs dropped the Labs from its name and now redis, the open source DB as well as Redis Labs, the commercial company, are referred to as "redis".
[Wikipedia](https://en.wikipedia.org/wiki/Redis)

#### Benefits

- Already available on Sqsc
- Available on AWS with Amazon ElastiCache for Redis
- Most popular DB
- Fully managed caching service that accelerates data access from primary databases and data stores with microsecond latency
- Supports 100,000 reads/writes per second
- 512MB for key, 512MB for value

#### Drawbacks

- Many different unofficial libraries for a simple message queue
- Basically redis is a key-value data storage system in memory not a message agent

## Technical resources

- [Cloudamqp] https://www.cloudamqp.com/blog/part1-rabbitmq-for-beginners-what-is-rabbitmq.html
