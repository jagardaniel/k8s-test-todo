# k8s-test-todo

This is more of a test repository to become familiar with with basic programming, Docker, Kubernetes and the build/deploy process. The goal is to create a simple todo application (frontend) and an API (backend) that it communicates with. GitHub Actions will then build containers on push and deploy them to Docker Hub. And after that we can hopefully push them to a Kubernetes cluster.


### Build and run in your local environment

Both the API and the todo application requires a LTS version of [Node.js](https://nodejs.org/en/) to build and run.

Clone the repository
```bash
$ git clone https://github.com/jagardaniel/k8s-test-todo.git
$ cd k8s-test-todo/
```

#### Todo API
The API is using the Express framework for Node.js. Most of it is based on this [gist](https://gist.github.com/colinskow/30ce0bf290db9b642ca456a15342f788). It stores everything in memory for now.

Enter the directory and install requirements with npm
```bash
$ cd todo-api/
todo-api$ npm install
```

You should then be able to run it
```bash
todo-api$ node src/index.js
```

...and reach it
```bash
$ curl localhost:5000/api/tasks
[{"id":1,"title":"Eat food","completed":true},{"id":2,"title":"Install something","completed":false},{"id":3,"title":"Eat food again","completed":true}]
```

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