scan/code: ## scans code for vulnerabilities
	@docker-compose --project-name trivy -f docker-compose.trivy.yml run --rm trivy fs /ingress-traefik-converter

scan/binary: ## scans binary for vulnerabilities
	@docker-compose --project-name trivy -f docker-compose.trivy.yml run --rm trivy fs /ingress-traefik-converter/dist/ingress-traefik-converter_darwin_amd64_v1/ingress-traefik-converter
