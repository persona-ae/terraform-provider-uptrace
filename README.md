# terraform-provider-uptrace

> This repo draws on the structure laid out by [this scaffolding repo](https://github.com/hashicorp/terraform-provider-scaffolding-framework)


## Pushing a new release:
> from [this article](https://thekevinwang.com/2023/10/05/build-publish-terraform-provider#github-release)

```bash
$ make release TAG=v0.0.3
```
goreleaser should create a new GitHub release with various artifacts included.

The Terraform Registry should detect this new release and create a new version — like magic.

That’s it!

> NOTE: if you see the error `template: failed to apply "{{ .Env.GPG_FINGERPRINT }}": map has no entry for key "GPG_FINGERPRINT"`
you'll need to export this like environment variable. See the GPG section below.


## Setting up Auth to Terraform Registry via GPG
To push from the cli you'll need your terraform registry account to have a gpg key connected. Generate this by running
`gpg --full-generate-key` and select the RSA and RSA option when prompted.
 When generating the key, enter 4096 for the key size and press Enter to accept the default option for the key expiration.
 Confirm your USER-ID by entering O. Add this key to the signing keys view. `$ gpg --export --armor "youremail@persona-ai.ai"`

Now export the environment variable to be consumed by the release script:
```bash
 gpg --list-keys

# /home/vscode/.gnupg/pubring.kbx
# -------------------------------
# pub   rsa4096 2021-12-15 [SC]
#       [[AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJ]]
# uid           [ultimate] Persona Dev <email@persona-ai.ai>
# sub   rsa4096 2021-12-15 [E]

export GPG_FINGERPRINT=[[AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJ]]
```
