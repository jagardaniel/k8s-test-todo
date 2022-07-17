function TaskItem(props) {
  return (
    <li className="list-group-item">
      <input
        id={"task-" + props.id}
        className="form-check-input me-3"
        type="checkbox" value=""
        defaultChecked={props.completed}
        onChange={() => props.toggleTaskCompleted(props.id, props.title, props.completed)}
      />
      {props.title}
      <button
        type="button"
        className="btn btn-danger badge bg-danger btn-sm float-end"
        onClick={() => props.deleteTask(props.id)}
      >
        X
      </button>
    </li>
  );
}

export default TaskItem;
