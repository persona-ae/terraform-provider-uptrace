# terraform-provider-uptrace

> This repo draws on the structure laid out by [this scaffolding repo](https://github.com/hashicorp/terraform-provider-scaffolding-framework)


## Pushing a new release:

You just need to make and push a new git tag and [our release action](./.github/workflows/release.yml) will do the rest!

```bash
git checkout main
git pull
git tag v0.0.8
git push origin v0.0.8
```
