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