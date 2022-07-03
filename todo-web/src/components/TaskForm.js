import React, { useState } from "react";

function TaskForm(props) {
  const [title, setTitle] = useState("");

  function handleSubmit(e) {
    e.preventDefault();
    props.addTask(title);
    setTitle("");
  }

  function handleChange(e) {
    setTitle(e.target.value);
  }

  return (
    <form onSubmit={handleSubmit}>
      <div className="mb-3">
        <input
          type="text"
          className="form-control"
          placeholder="Create new task"
          id="inputTask"
          value={title}
          onChange={handleChange}
        />
      </div>
      <button type="submit" className="btn btn-primary">Submit</button>
    </form>
  )
}

export default TaskForm;
