ACTION_RUNNER_ARCH = x64
ACTION_RUNNER_VER = 2.299.1
ACTION_RUNNER_ARCHIVE = actions-runner-linux-$(ACTION_RUNNER_ARCH)-$(ACTION_RUNNER_VER).tar.gz

gh-action-runner:	## Install Github action custom runner.
	mkdir -p actions-runner && cd actions-runner
	curl -O -L https://github.com/actions/runner/releases/download\
	/v$(ACTION_RUNNER_VER)/actions-runner-linux-$(ACTION_RUNNER_ARCH)\
	-$(ACTION_RUNNER_VER).tar.gz
	tar xzf ./$(ACTION_RUNNER_ARCHIVE)
	rm $(ACTION_RUNNER_ARCHIVE)
	./config.sh --url $(GITHUB_REPO) --token $(ACTION_RUNNER_TOKEN)

gh-action-run:		## Run the Github action custom runner.
	./run.sh
