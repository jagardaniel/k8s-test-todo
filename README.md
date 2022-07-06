# k8s-test-todo

This is more of a test repository to become familiar with with basic programming, Docker, Kubernetes and the build/deploy process. The goal is to create a simple todo application (frontend) and an API (backend) that it communicates with. GitHub Actions will then build containers on push and deploy them to Docker Hub. And after that we can hopefully push them to a Kubernetes cluster.


### Build and run in your local environment

Clone the repository
```bash
$ git clone https://github.com/jagardaniel/k8s-test-todo.git
$ cd k8s-test-todo/
```

#### Todo API

#### Todo Web
The frontend is written in React. Mozillas React [tutorial](https://developer.mozilla.org/en-US/docs/Learn/Tools_and_testing/Client-side_JavaScript_frameworks/React_getting_started) and [this](https://github.com/imnileshd/react-todo-app) repository has been two great resources for my desperate attempt to create something that actually starts.


Enter the directory and install requirements with npm
```bash
$ cd todo-web/
todo-web$ npm install
```

Run `npm start` to start the application. It will open up a new tab with the application in your browser automatically.
```bash
todo-web$ npm start
```

If I haven't done something stupid (it is possible!) and the API is running you should now be able to list, create, toggle and delete tasks. The development server is configured to proxy requests to the API, so a request to `/api/tasks` for example will be sent to `http://localhost:5000/api/tasks`. This means that we should be able to avoid CORS issues.

#### Docker
Soon!