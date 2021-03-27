ACTION_RUNNER_ARCH = x64
ACTION_RUNNER_VER = 2.277.1

gh-action-runner:	## Install Github action custom runner.
	mkdir -p actions-runner && cd actions-runner
	curl -O -L https://github.com/actions/runner/releases/download\
	/v$(ACTION_RUNNER_VER)/actions-runner-linux-$(ACTION_RUNNER_ARCH)\
	-$(ACTION_RUNNER_VER).tar.gz
	tar xzf ./actions-runner-linux-$(ACTION_RUNNER_ARCH)-$(ACTION_RUNNER_VER).tar.gz
	./config.sh --url $($(GITHUB_REPO)) --token $(ACTION_RUNNER_TOKEN)

gh-action-run:		## Run the Github action custom runner.
	./run.sh
