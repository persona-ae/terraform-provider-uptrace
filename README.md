# terraform-provider-uptrace

> This repo draws on the structure laid out by [this scaffolding repo](https://github.com/hashicorp/terraform-provider-scaffolding-framework)


## Pushing a new release:
> from [this article](https://thekevinwang.com/2023/10/05/build-publish-terraform-provider#github-release)

```bash
git tag [[v0.1.1]]
git push origin [[v0.1.1]]
GITHUB_TOKEN=$(gh auth token) goreleaser release --clean
```
goreleaser should create a new GitHub release with various artifacts included.

The Terraform Registry should detect this new release and create a new version — like magic.

That’s it!


