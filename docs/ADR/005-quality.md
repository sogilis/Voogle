# ADR-005-quality

- Creation Date: 14/01/2022
- Status: Accepted

## Tools

### Continuous integration

Continuous integration is the ability of a software factory to automatically and continuously verify the proper construction (build), validation (automated tests) and even delivery of the software.

In other words, continuous integration is the automatic process of validating and generating artifacts, making them available and deploying them to a given environment. This allows for the continuous verification of the software's ability to be delivered and deployed. It also allows to always have the latest version produced by the developers.

By using continuous integration, it is easy to split the production of different versions of the software (beta, release candidate or stable), which facilitates the involvement of transverse teams (for example to validate the transition from beta to release candidate).

### Coding standard and code ergonomics

#### Basic principle

The Sogilis team adopts coding standards inspired by the Clean Code reference(https://dl.acm.org/doi/book/10.5555/1388398).
These coding standards allow us to produce good code ergonomics in the sense of the "Bastien and Scapin criteria". (https://hal.inria.fr/inria-00070476v2/document).

The coding standards adopted by Sogilis define coding rules for the following elements: guidance, explicit control, workload, adaptability, error management, homogeneity/consistency, meaning of codes and names and compatibility.

#### Golang specific style

For golang specific style standards, the following references will be used, in order of priority:

1. https://golang.org/doc/effective_go
2. https://github.com/uber-go/guide/blob/master/style.md
3. https://docs.gitlab.com/ee/development/go_guide/

#### Linter

A linter is used to verify and ensure (its launch in a validation stage of the CI) the good formatting of the Go code of the project.
The linter used at the moment is [golangci-lint](https://github.com/golangci/golangci-lint).
