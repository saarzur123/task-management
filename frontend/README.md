# Overview

A ReactJS application that integrates with the backend, enabling users to perform CRUD operations on tasks through an intuitive user interface.

## Components
### TaskActionsModal.js
`TaskActionsModal.js` is a React component responsible for managing the creation and editing of tasks through a modal dialog.  
It enables users to input or modify task details, such as the title, description, and status, with required fields enforced.  
The modal dynamically adapts based on whether the user is creating a new task or updating an existing one.  
It interacts with the backend API to perform task creation or updates, while displaying loading indicators during the process and handling errors in case of failures.

### TasksTable.js
`TasksTable.js` is a React component that displays a table of tasks fetched from the backend API.  
It shows each task's title, description, and status, and provides actions for editing and deleting tasks.  
Users can also add new tasks through an "Add" button (+ icon).  
The component manages the state for task data, handling loading states and dynamically updating the task list after actions such as creating, editing, or deleting tasks.  
It ensures smooth integration with the backend and provides a clean, organized table view for task management.

## Technology Stack
The frontend of the application is built using **React.js** which makes it easy to build reusable UI elements and manage the application state.  
The state is managed using React’s built-in `useState` hook, which allows for tracking dynamic values like task data and loading states.  
The `useEffect` hook is used to handle side effects, such as fetching data from the backend API when the component mounts or updates.  
For API communication, the fetch API is used to make HTTP requests, allowing the app to interact with the backend for CRUD operations on tasks.

# Setup
### Prerequisites
- [Node.js](https://nodejs.org/) (optional for local frontend development)

### Run Locally
  - Navigate to the `frontend` folder.
  - Install dependencies:
    ```bash
    npm install
    ```
  - Start the development server:
    ```bash
    npm start
    ```