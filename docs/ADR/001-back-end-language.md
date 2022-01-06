# ADR- 001 - Back-end language

Creation Date: 05/01/2022

## Context

In order to write the first bricks of the Voogle application, we need to choose a language that suites our needs.
We already have identified some languages, NodeJS, Go and Java Quarkus.

The chosen language must be suited to write Cloud Native application.

## Decision

## Status

Draft

## Options

### 1. NodeJS

*"Node.js has an event-driven architecture capable of asynchronous I/O. These design choices aim to optimize throughput and scalability in web applications with many input/output operations"*
 [Wikipedia](https://en.wikipedia.org/wiki/Node.js)

#### Benefits
* NodeJS is one of the most popular language
* Rich Ecosystem, there are many package available and almost all the topics are covered
* NodeJS is rather performant
* JavaScript is a language that allow to develop quickly
* NodeJS's applications start quickly

#### Drawbacks
* Heavy treatment can impact performances
* If not used properly, it can the event programming paradigm can have bad consequences performance wise
* Requires the interpreter to be installed
* The Async system can be hard to maintain


#### Points of attention
* Typescript is superset of NodeJS that allows to have a strongly typed language on top of Javascript. It allows us to catch a lot of mistakes that can be 
  make with JavaScript. Plus Typescript is transpiled to native JavaScript.

### Go

#### Benefits
#### Drawbacks
#### Points of attention

### Java Quarkus
Quarkus is a Java Framework that was conceived for Cloud application. It aims to be lighter and quicker to start that standard Java application (especially 
when using the Native build).

#### Benefits
* Quick boot time (99% less than traditional Java app)
* Smaller memory footprint (86% less than traditional Java app)
* Can be compiled into a native application (Without JVM)
* Live reload is supported
* Reactive and imperative development paradigm (At the same time)

#### Drawbacks
* Much longer build (10x longer)
* With native build, some Java features (like reflections) are not available
* Some incompatibilities with some major libraries (like Spring boot (partially))
* Some pre-configurations already done by default 

#### Points of attention

| /            | Learning curve | Horizontal scalability | Vertical Scalability | CPU Bound treatment | Maintainability/Ease of deployment | Productivity |
|--------------|----------------|------------------------|----------------------|---------------------|------------------------------------|--------------|
| NodeJS       | ++             | +                      | +/-                  | -                   | -                                  | +            |
| Go           | ?              | ?                      | ?                    | ?                   | ?                                  | ?            |
| Java Quarkus | +-             | +                      | +                    | +                   | +                                  | +-           |

## Technical resources
* NodeJS and Heavy CPU bound task: http://neilk.net/blog/2013/04/30/why-you-should-use-nodejs-for-CPU-bound-tasks/
* NodeJS vs Go performance benchmark: https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go-node.html
* Morgan's Quarkus REX (FR): https://www.notion.so/Quarkus-Supersonic-Subatomic-Java-a89562f563c14d07a1afbe206f6f7da8
* Quarkus main page with some "benchmarks": https://quarkus.io/
* Wikipedia page about Quarkus (And the default configuration): https://en.wikipedia.org/wiki/Quarkus#Zero_configuration_with_dev_services


