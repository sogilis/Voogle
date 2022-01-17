# GitHub Container Registry - [ghcr.io](ghcr.io)

## How to do it locally
```bash
echo $PAT | docker login ghcr.io --username phanatic --password-stdin

docker tag app ghcr.io/phanatic/app:1.0.0

docker push ghcr.io/phanatic/app:1.0.0
```

## Links:

* [Working with the container registry](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry)
* [Security hardening for GitHub Actions ](https://docs.github.com/en/actions/security-guides/security-hardening-for-github-actions#considering-cross-repository-access)
* [Setup security parameters](https://docs.github.com/en/packages/managing-github-packages-using-github-actions-workflows/publishing-and-installing-a-package-with-github-actions#upgrading-a-workflow-that-accesses-ghcrio)
