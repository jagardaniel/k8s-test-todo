# k8s-test-todo

This is more of a test repository to become familiar with with basic programming, Docker, Kubernetes and the build/deploy process. The goal is to create a simple todo application (frontend) and an API (backend) that it communicates with. GitHub Actions will then build containers on push and deploy them to Docker Hub. And after that we can hopefully push them to a Kubernetes cluster.


### Build and run in your local environment

Clone the repository
```bash
$ git clone https://github.com/jagardaniel/k8s-test-todo.git
$ cd k8s-test-todo/
```

#### Todo API
The API is written in Go and the web framework Gin. I have mostly used [Learn Go With Tests](https://quii.gitbook.io/learn-go-with-tests/build-an-application/app-intro) and [this blog series](https://eli.thegreenplace.net/2021/rest-servers-in-go-part-1-standard-library/) as resources to build something that I can write tests for, since that is a part I want to include in the CI workflow. So the parts of the code that actually looks good are probably from one of these two great resources.

To run the API you need a somewhat recent version of [Go](https://go.dev/dl/). Then build and run it.
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
