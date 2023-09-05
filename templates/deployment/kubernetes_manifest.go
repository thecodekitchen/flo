package deployment

func KubernetesManifestBytes(project string, registry string) []byte {
	return []byte(
		`apiVersion: apps/v1
kind: Deployment 
metadata:
  name: api
spec:
  selector:
	matchLabels:
	  app: api
  template:
	metadata:
	  labels:
		app: api
	  spec:
	    containers:
        - name: api
          image: ` + registry + `/` + project + `
          ports:
          - containerPort: 8080
          env:
          - name: SURREALDB_URL
            value: ws://surrealdb-tikv:8000/rpc
---

# Public API Service
apiVersion: v1
kind: Service
metadata:
  name: api-service 
spec:
  selector:
	app: api
  type: LoadBalancer  
  ports:
	- protocol: TCP
	  port: 8080
	  targetPort: 8080`)
}
