template from: https://git.matador.ais.co.th/cronus/irm/devops/gitlab-init-template.git


device-register
741W36Cz&jtu

helm template ./helm -f helm/values.yaml --name-template mobile-validator-register-service --set deployment.image="registry.gitlab.com/nextdb-project/digital-reality-foundation/mobile-validator/validator-register/mobile-validator-scheduler-service:latest" -n transaction-gateway-be > helm/k8s/deployment.yml


helm install ./helm -f helm/values.yaml --name-template mobile-validator-register-service --set deployment.image="registry.gitlab.com/nextdb-project/digital-reality-foundation/mobile-validator/validator-register/mobile-validator-scheduler-service:latest" -n transaction-gateway-be

kubectl create secret docker-registry gitlab-transaction-gateway-cr --docker-server=registry.gitlab.com --docker-username=gitlab+deploy-token-1256371 --docker-password=dwvr7wxH11CEZ62bLhGC --docker-email=wasawatp@ais.co.th