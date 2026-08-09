all: image-build image-push

image-build:
	podman build -t git.servidordomal.lol/robogg133/superdeploy:latest .

image-push:
	podman push git.servidordomal.lol/robogg133/superdeploy:latest

rollout:
	kubectl -n superdeploy-namespace rollout restart deployment/superdeploy-app
