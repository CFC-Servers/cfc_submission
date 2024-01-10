build-and-deploy:
	/usr/local/bin/sam build && /usr/local/bin/sam deploy --no-confirm-changeset --no-fail-on-empty-changeset
