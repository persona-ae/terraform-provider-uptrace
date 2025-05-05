.PHONY: ensure-go-releaser
ensure-go-releaser:
	@which goreleaser >/dev/null 2>&1 || (\
		echo "goreleaser not found. installing..."; \
		echo 'deb [trusted=yes] https://repo.goreleaser.com/apt/ /' | sudo tee /etc/apt/sources.list.d/goreleaser.list; \
		sudo apt update; \
		sudo apt install goreleaser; \
	)

.PHONY: ensure-gh
ensure-gh:
	@which gh >/dev/null 2>&1 || (\
		echo "gh not found. installing..."; \
		(type -p wget >/dev/null || (sudo apt update && sudo apt-get install wget -y)) \
		&& sudo mkdir -p -m 755 /etc/apt/keyrings \
		&& out=$(mktemp) && wget -nv -O$out https://cli.github.com/packages/githubcli-archive-keyring.gpg \
		&& cat $out | sudo tee /etc/apt/keyrings/githubcli-archive-keyring.gpg > /dev/null \
		&& sudo chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg \
		&& echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null \
		&& sudo apt update \
		&& sudo apt install gh -y; \
	)

release: ensure-go-releaser ensure-gh
	@GITHUB_TOKEN=$(gh auth token) goreleaser release --clean
