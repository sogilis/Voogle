# ADR-006-Videos-storage

* Creation Date: 20/02/2022
* Status: Draft 

## Context
Storing our videos on the API filesystem is painful and not scalable. We can't easily add new videos and if you want multiple instances of API, they all need to have copies of the videos.

## Decision

We will use a S3 storage for our videos. For the production environment, we will use Amazon S3 and for development purpose, we will use MinIO.

## Options
### Use S3 storage
#### Benefits
* Highly scalable
* Can be deployed with terraform
* Will allow us to use terraform and S3
* Easily queried
#### Drawbacks
* Slower to setup

## Technical resources
- [AWS S3](https://aws.amazon.com/fr/s3/)
- [MinIO](https://min.io/)
