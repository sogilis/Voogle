# ADR-003-Continuous-Integration

* Creation Date: 06/01/2022
* Status: Draft

## Context

With our goal of quality and that we want to also do Continuous Deployment, we need a tool do execute our tests, build our artifacts and deploy them.
We have identified a few tools (ConcourseCI, GitHub Actions, DroneCI and TeamCity), and we need to pick one.

## Decision

ToDo

## Options

### 1.ConcourseCI

#### Benefits
* Open source
* Nice and simple interface

#### Drawbacks
* ConcourseCI is only available on premise (Some tiers provide a SaaS version)
* Young
* Integration with external source is not intuitive (Documentation is sparse)
* YAML for pipeline configuration
* I (JPR) was a bit disappointed by the tool when I tried it, the tool doesn't appear intuitive

### 2. GitHub Actions

#### Benefits
* We already use GitHub as our source control tool, so the integration is effortless
* We already pay for GitHub Pro, it might not increase our costs
* A community market of plugins
* [Community Tool](https://github.com/nektos/act) to run the pipelines locally
* We already have some experience with it

#### Drawbacks
* Closed source
* Not Free
* Use YAML for the configuration

### 3. DroneCI

*"Drone is a self-service Continuous Integration platform for busy development teams."*
[Drone.io](https://www.drone.io/)
#### Benefits
* Open source for open source project 
* Designed for docker
* Easy integration with GitHub
* Syntax is rather easy to read and quite explicit, configuration syntax is a derivative of docker-compose
* Many plugins (docker container)

#### Drawbacks
* Configuration with YAML
* On premise or Enterprise Edition $299/month for 25 Developers

### 4. TeamCity

#### Benefits
* Configuration file in Kotlin
* Easy integration with GitHub with webhooks
* Code Quality tools integrated
* Integrated with JetBrains tools
#### Drawbacks
* On premise or 45$ per month

## Technical resources
- [GitHub Actions features](https://github.com/features/actions)
- [Article about complex pipeline in GitHub Actions](https://dh1tw.de/2019/12/real-life-ci/cd-pipelines-with-github-actions/)
- [Example of GitHub Actions workflow](https://github.com/dh1tw/remoteAudio/blob/master/.github/workflows/build.yml)
- [TeamCity](https://www.jetbrains.com/teamcity/)
- [TeamCity & GitHub](https://ardalis.com/4-tips-to-integrate-teamcity-and-github/)
- [Drone-ci drone cloud](https://blog.drone.io/drone-cloud/)
- [Drone-ci plugins](http://plugins.drone.io/)