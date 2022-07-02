const express = require("express")

const app = express()
const port = 3000

app.use(express.json())

// Add some default tasks
let tasks = [
  {id: 1, title: "Eat food", completed: true},
  {id: 2, title: "Install something", completed: false},
  {id: 3, title: "Eat food again", completed: true},
]
let nextId = 4;

// GET: /api/tasks
app.get("/api/tasks", (req, res) => {
  res.status(200).json(tasks)
})

// POST: /api/tasks
app.post("/api/tasks", (req, res, next) => {
  if (!req.body || !req.body.title) {
    return next({status: 400, message: "title field is required"})
  }

  if (req.body.title.length < 3) {
    return next({status: 400, message: "title has to be at least 3 characters long"})
  }

  const newTask = {
    id: nextId,
    title: req.body.title,
    completed: false
  }
  tasks.push(newTask)

  nextId++

  res.status(201).json(newTask)
})

// GET: /api/tasks/:id
app.get("/api/tasks/:id", (req, res, next) => {
  const id = parseInt(req.params.id, 10)

  const task = tasks.find((item) => item.id === id)
  if (!task) {
    return next({status: 404, message: "task not found"})
  }

  res.status(200).json(task)
})

// DELETE: /api/tasks/:id
app.delete("/api/tasks/:id", (req, res) => {
  const id = parseInt(req.params.id, 10)

  const taskIndex = tasks.findIndex((item) => item.id === id)
  if (taskIndex !== -1) {
    tasks.splice(taskIndex, 1)
  }
  
  res.sendStatus(204)
})

// Toggle "completed" value to true
// POST: /api/tasks/:id/completed
app.post("/api/tasks/:id/completed", (req, res, next) => {
  const id = parseInt(req.params.id, 10)

  const taskIndex = tasks.findIndex((item) => item.id === id)
  if (taskIndex === -1) {
    return next({status: 404, message: "task not found"})
  }
  
  tasks[taskIndex].completed = true

  res.sendStatus(201)
})

// Toggle "completed" value to false
// DELETE: /api/tasks/:id/completed
app.delete("/api/tasks/:id/completed", (req, res, next) => {
  const id = parseInt(req.params.id, 10)

  const taskIndex = tasks.findIndex((item) => item.id === id)
  if (taskIndex === -1) {
    return next({status: 404, message: "task not found"})
  }
  
  tasks[taskIndex].completed = false

  res.sendStatus(204)
})

// Error handler
app.use((err, req, res, next) => {
  res.status(err.status).json(err)
})

app.listen(port, () => {
  console.log(`Listening on port ${port}`)
})
