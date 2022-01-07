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
* 

#### Drawbacks
* ConcourseCI is only available on premise (Some tiers provide a SaaS version)
* Young
* Integration with external source is not intuitive
* YAML for pipeline configuration

#### Points of attention

### 2. GitHub Actions

#### Benefits
* We already use GitHub as our source control tool, so the integration is effortless
* We already pay for GitHub Pro, it might not increase our costs
* A community market of plugins

#### Drawbacks
* Closed source
* Not Free
* Use YAML for the configuration

### 3. DroneCI

#### Benefits
#### Drawbacks
#### Points of attention

### 4. TeamCity

#### Benefits
#### Drawbacks
#### Points of attention

## Technical resources
- [GitHub Actions features](https://github.com/features/actions)
- [Article about complex pipeline in GitHub Actions](https://dh1tw.de/2019/12/real-life-ci/cd-pipelines-with-github-actions/)
- [Example of GitHub Actions workflow](https://github.com/dh1tw/remoteAudio/blob/master/.github/workflows/build.yml)
