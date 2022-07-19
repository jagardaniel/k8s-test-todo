# k8s-test-todo

This is a test repository to become more familiar with with programming, Docker, Kubernetes and the build/deploy process. The goal is to create a simple todo web application and an API that it communicates with to store and retrieve todo tasks. GitHub Actions will then build/test, build containers and deploy them to Docker Hub. After that we can hopefully push them to a Kubernetes cluster.

The difference from the previous k8s-test repository is that the frontend is written in React and performs the API calls on the client side. This makes it a little bit tricker since we have to work around CORS issues.


### Build and run in your local environment

Clone the repository
```bash
$ git clone https://github.com/jagardaniel/k8s-test-todo.git
$ cd k8s-test-todo/
```

#### Todo API
The API is written in Go and the web framework Gin. I have mostly used [Learn Go With Tests](https://quii.gitbook.io/learn-go-with-tests/build-an-application/app-intro) and [this blog series](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/) as resources to build something that I can write tests for, since that is a part I want to include in the CI workflow. So the parts of the code that actually looks good are probably from one of these two great resources.

To run the API you need a somewhat recent version of [Go](https://go.dev/dl/). Then build and run the binary.
```bash
$ cd todo-api/
todo-api$ go build
todo-api$ ./todo-api
```

You should now be able to reach it on port `:8000`.

```bash
todo-api$ curl 127.0.0.1:8000/api/tasks
[]
```

#### Todo Web
The frontend is written in React. Mozillas React [tutorial](https://developer.mozilla.org/en-US/docs/Learn/Tools_and_testing/Client-side_JavaScript_frameworks/React_getting_started) and [this](https://github.com/imnileshd/react-todo-app) repository has been two great resources for my desperate attempt to create something that actually starts. It requires a recent version of [Node.js](https://nodejs.org/en/).


Enter the directory and install requirements with npm
```bash
$ cd todo-web/
todo-web$ npm install
```

Run `npm start` to start the application. It will open up a new tab with the application in your browser automatically.
```bash
todo-web$ npm start
```

If I haven't done something stupid (it is possible!) and the API is running you should now be able to list, create and delete tasks. The development server is configured to proxy requests to the API, so a request to `/api/tasks` will be sent to `http://localhost:8000/api/tasks` for example. This means that we should be able to avoid CORS issues.

### CI/CD

### Kubernetes

Like the previous time I don't have access to any fancy Kubernetes clusters with cool names so minikube will do it for us this time as well. We need to enable the ingress addon for the ingress part to work. We also need to point the ingress host to the ingress IP address. I use a non-existing subdomain (`k8s-test.bottenskrap.se`) in a zone that I own for this since we run everything in a local test environment.

Start minikube
```bash
$ minikube start
```

Enable the ingress addon and verify that it is running
```bash
# Enable addon
$ minikube addons enable ingress

# Verify (should look something like this)
$ kubectl get pods -n ingress-nginx
NAME                                       READY   STATUS      RESTARTS   AGE
ingress-nginx-admission-create-s46fl       0/1     Completed   0          35s
ingress-nginx-admission-patch-rt9jc        0/1     Completed   0          35s
ingress-nginx-controller-cc8496874-xjpjr   1/1     Running     0          35s
```

Deploy deployments, services and ingress
```bash
$ kubectl apply -f k8s/
ingress.networking.k8s.io/todo-ingress created
deployment.apps/todo-api created
service/todo-api created
deployment.apps/todo-web created
service/todo-web created
```

Get the ingress IP address so we can create a host entry for it
```bash
$ kubectl get ingress
NAME           CLASS   HOSTS                     ADDRESS          PORTS   AGE
todo-ingress   nginx   k8s-test.bottenskrap.se   192.168.39.167   80      60s
```

Add the host entry. Again, I just use a domain that does not exist. If you want to use another domain make sure to match the host value in the `ingress.yaml` file.
```bash
echo "192.168.39.167 k8s-test.bottenskrap.se" | sudo tee -a /etc/hosts
```

You should now be able to visit the domain in your browser and use the todo application!

Then a couple of seconds later when you realize how unnecessary the application is you can run `kubectl delete` to remove everything you just deployed from the cluster
```bash
$ kubectl delete -f k8s/
ingress.networking.k8s.io "todo-ingress" deleted
deployment.apps "todo-api" deleted
service "todo-api" deleted
deployment.apps "todo-web" deleted
service "todo-web" deleted
```

### Red Hat OpenShift

Red Hat provides a 30-day sandbox environment for OpenShift which is great for experimenting and learning Kubernetes and OpenShift. More information [here](https://developers.redhat.com/developer-sandbox/activities). Lets try and push our two containers using their CLI tool. I don't know if this is the "recommended" way to deploy applications in OpenShift.

After we have created a sandbox environment we need to install the `oc` CLI tool. Instructions [here](https://docs.openshift.com/container-platform/4.7/cli_reference/openshift_cli/getting-started-cli.html).


Use `oc login` to login with your token. You can find the complete login command if you click on your username in the right top corner from the web interface and then select `Copy login command`.
```bash
$ oc login --token=YOURTOKENHERE --server=https://api.sandbox.xxxx.xx.openshiftapps.com:6443
Logged into "https://api.sandbox.xxxx.xx.openshiftapps.com:6443" as "jagardaniel" using the token provided.

You have one project on this server: "jagardaniel-dev"

Using project "jagardaniel-dev".
```

We can use the command `oc new-app image-repo-url` to create both our applications. The tool is clever and can figure out name, language, port and other things based on the image. It then creates a deployment, service and an ImageStream(?) in OpenShift for us. If you run the command with `-o yaml` you can see the definition without creating the application.

```bash
# Deploy todo-api
$ oc new-app docker.io/jagardaniel/todo-api
--> Found container image 81d61ca (3 hours old) from docker.io for "docker.io/jagardaniel/todo-api"
[...]

# Deploy todo-web
$ oc new-app docker.io/jagardaniel/todo-web
--> Found container image f54aa95 (3 hours old) from docker.io for "docker.io/jagardaniel/todo-web"
[...]

# Verify
$ oc get deployments,pods,services
NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/todo-api   1/1     1            1           11m
deployment.apps/todo-web   1/1     1            1           6m43s

NAME                            READY   STATUS    RESTARTS   AGE
pod/todo-api-7c5d759756-7b4c8   1/1     Running   0          11m
pod/todo-web-864ddf86d9-rqkwv   1/1     Running   0          6m43s

NAME               TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
service/todo-api   ClusterIP   172.30.188.46    <none>        8000/TCP   11m
service/todo-web   ClusterIP   172.30.174.250   <none>        8080/TCP   6m43s
```

We then need to create two routes so we can reach them externally.

```bash
# Create route for the todo-api service
$ oc expose --hostname todo-test.bottenskrap.se --path /api service/todo-api
route.route.openshift.io/todo-api exposed

# Create route for the todo-web service
$ oc expose --hostname todo-test.bottenskrap.se service/todo-web
route.route.openshift.io/todo-web exposed

# Verify
$ oc get routes
NAME       HOST/PORT                  PATH   SERVICES   PORT       TERMINATION   WILDCARD
todo-api   todo-test.bottenskrap.se   /api   todo-api   8000-tcp                 None
todo-web   todo-test.bottenskrap.se          todo-web   8080-tcp                 None
```

We also need to create a CNAME record that points our specified hostname (`todo-test.bottenskrap.se` in this example) to the "Router canonical hostname". You can use `oc describe` to find it.

```bash
$ oc describe route/todo-api | grep -A 1 Requested
Requested Host:		todo-test.bottenskrap.se
			   exposed on router default (host router-default.apps.sandbox.xxxx.xx.openshiftapps.com) 39 minutes ago
```

So in the example above it would be `router-default.apps.sandbox.xxxx.xx.openshiftapps.com)`. So a query should look something like this:

```bash
$ dig +short todo-test.bottenskrap.se
router-default.apps.sandbox.xxxx.xx.openshiftapps.com.
[...]
```

And you should now be able to visit the hostname in a browser!

When you get bored you can remove them with `oc delete`.
```bash
$ oc delete all --selector app=todo-web
service "todo-web" deleted
deployment.apps "todo-web" deleted
imagestream.image.openshift.io "todo-web" deleted
route.route.openshift.io "todo-web" deleted

$ oc delete all --selector app=todo-api
service "todo-api" deleted
deployment.apps "todo-api" deleted
imagestream.image.openshift.io "todo-api" deleted
route.route.openshift.io "todo-api" deleted
```