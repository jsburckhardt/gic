configure-commit-signature:
	@echo "Signing commit..."
	@bash ./src/scripts/sign_commits.sh jsburckhardt jsburckhardt@gmail.com

purge-resources:
	@echo "Purging resources..."
	@bash ./src/scripts/purge_resources.sh
