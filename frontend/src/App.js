import './App.css';
import TasksTable from "./components/TasksTable";

function App() {
  return (
    <div className="tasks-app">
        <h1>Tasks Manager</h1>
      <TasksTable/>
    </div>
  );
}

export default App;
