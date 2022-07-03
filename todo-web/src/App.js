import React, { useState, useEffect } from "react";
import TaskItem from "./components/TaskItem";
import TaskForm from "./components/TaskForm";

function App() {
  const [error, setError] = useState(null);
  const [isLoaded, setIsLoaded] = useState(false);
  const [tasks, setTasks] = useState([]);

  // Code to run on component mount
  useEffect(() => {
    getTasks();
  }, [])

  function getTasks() {
    fetch("/api/tasks")
      .then(res => res.json())
      .then(
        (result) => {
          setIsLoaded(true);
          setTasks(result);
        },
        (error) => {
          setIsLoaded(true);
          setError(error);
        }
      )
  }

  function addTask(title) {
    const requestOptions = {
      method: "POST",
      headers: {"Content-Type": "application/json"},
      body: JSON.stringify({title: title})
    };

    fetch("/api/tasks", requestOptions)
      .then(res => res.json())
      .then(
        (result) => {
          getTasks();
        },
        (error) => {
          console.log(error);
        }
      )
  }

  function deleteTask(id) {
    fetch("/api/tasks/" + id, {method: "DELETE"})
      .then(
        () => getTasks()
      );
  }  

  function toggleTaskCompleted(id, completed) {
    if (completed) {
      fetch("/api/tasks/" + id + "/completed", {method: "DELETE"})
    } else {
      fetch("/api/tasks/" + id + "/completed", {method: "POST"})
    }
  }

  function renderTaskList() {
    if (error) {
      return (
          <div className="alert alert-danger" role="alert">
            Could not fetch tasks. Is the API up and running?
          </div>
        );
    } else if (!isLoaded) {
      return (
        <div className="d-flex justify-content-center">
          <div className="spinner-border text-primary" role="status">
            <span className="visually-hidden">Loading...</span>
          </div>
        </div>
      );
    } else {
      if (tasks.length > 0) {
        return <ul className="list-group">{taskList}</ul>;
      } else {
        return <p className="card-text">Nothing to do!</p>
      }
    }
  }

  const taskList = tasks.map(task => (
    <TaskItem
      id={task.id}
      title={task.title}
      completed={task.completed}
      key={"task-" + task.id}
      toggleTaskCompleted={toggleTaskCompleted}
      deleteTask={deleteTask}
    />
  ));

  return (
    <div className="container">
      <div className="p-3 bg-body rounded shadow-sm">
        <h2>Todo List</h2>
        <TaskForm addTask={addTask} />
        <hr />
        {renderTaskList()}
      </div>
    </div>
  );
}

export default App;
