deps-darwin:
	@brew cask install xquartz
	@brew install socat

get-addons:
	# Skip smudge - We'll download binary files later in a faster batch
	# @git lfs install --skip-smudge
	# Do git clone here
	@chmod a+x ./shared/scripts/subtrees.sh
	@./shared/scripts/subtrees.sh
	# Fetch all the binary files in the new clone
	#@export GIT_SSL_NO_VERIFY=1 && git lfs pull
	# Reinstate smudge
	#@git lfs install --force