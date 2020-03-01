### K8s Disrupter Server 🔥

```
kubectl create namespace gremlin

kubectl create secret generic gremlin-team-cert \
 --namespace=gremlin \
 --from-file=gremlin.cert \
 --from-file=gremlin.key


GREMLIN_TEAM_ID="GREMLIN_TEAM_ID"
GREMLIN_CLUSTER_ID="chaotic-cluster"

helm repo remove gremlin
helm repo add gremlin https://helm.gremlin.com

helm install \
 --namespace gremlin gremlin \
 gremlin/gremlin \
 --set gremlin.teamID=$GREMLIN_TEAM_ID \
	--set gremlin.clusterID=$GREMLIN_CLUSTER_ID

gcloud app deploy
```
