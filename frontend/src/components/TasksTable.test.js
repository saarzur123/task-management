import React from "react";
import { act } from 'react';
import {render, screen, fireEvent, waitFor, cleanup} from "@testing-library/react";
import TasksTable from "./TasksTable";

global.fetch = jest.fn();
global.alert = jest.fn();

describe("TasksTable", () => {
    const baseTask = {
        id: 1,
        title: "Task 1",
        description: "Test Task",
        status: "Pending",
    };
    const mockTasks = [
        baseTask,
        { id: 2, title: "Task 2", description: "Test Task2", status: "Pending2" },
    ];
    const mockNewTask = {
        id: 2,
        title: "New Task",
        description: "New description",
        status: "In Progress",
    };
    const errorMessage = "Failed to fetch tasks";
    const noTasksMessage = "No tasks to show, use the create (+) button to add tasks.";

    beforeEach(() => {
        global.alert.mockClear();
    });

    afterEach(() => {
        jest.clearAllMocks();
        cleanup();
    });

    test("renders tasks when fetched successfully", async () => {
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => mockTasks,
        });
        render(<TasksTable/>);
        await waitFor(() => screen.getByText("Task 1"));
        await waitFor(() => screen.getByText("Task 2"));

        expect(screen.getByText("Task 1")).toBeInTheDocument();
        expect(screen.getByText("Test Task")).toBeInTheDocument();
        expect(screen.getByText("Pending")).toBeInTheDocument();
        expect(screen.getByText("Task 2")).toBeInTheDocument();
        expect(screen.getByText("Test Task2")).toBeInTheDocument();
        expect(screen.getByText("Pending2")).toBeInTheDocument();
    });

    test('renders error message if fetch fails', async () => {
        fetch.mockRejectedValueOnce(new Error(errorMessage));
        render(<TasksTable/>)
        await waitFor(() => expect(global.alert).toHaveBeenCalledWith(errorMessage));

        expect(global.alert).toHaveBeenCalledTimes(1);
        expect(global.alert).toHaveBeenCalledWith(errorMessage);
    });

    test("opens the TaskActionsModal when the add task icon is clicked", async () => {
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => [baseTask],
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));

        fireEvent.click(screen.getByTestId("AddCircleIcon"));
        expect(screen.getByLabelText("Create Task")).toBeInTheDocument();
        expect(screen.queryByText(noTasksMessage)).not.toBeInTheDocument();
    });

    test("opens the TaskActionsModal when the edit icon is clicked", async () => {
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => [baseTask],
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));

        fireEvent.click(screen.getByTestId("EditIcon"));
        expect(screen.getByLabelText(/Edit Task 1/i)).toBeInTheDocument();
    });

    test("calls deleteTask and removes specific task from the list", async () => {
       // add data to the table before deletion
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => mockTasks,
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));
        await waitFor(() => screen.getByText("Task 2"));

        // delete
        fetch.mockResolvedValueOnce({ ok: true });
        fireEvent.click(screen.getAllByTestId("DeleteIcon")[0]);

        await waitFor(() => expect(screen.queryByText("Task 1")).not.toBeInTheDocument());
        expect(screen.getByText("Task 2")).toBeInTheDocument();
        expect(screen.queryByText(noTasksMessage)).not.toBeInTheDocument();
    });

    test("calls deleteTask and removes all the tasks from the list", async () => {
        // add data to the table before deletion
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => mockTasks,
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));
        await waitFor(() => screen.getByText("Task 2"));

        // delete first task
        fetch.mockResolvedValueOnce({ ok: true });
        fireEvent.click(screen.getAllByTestId("DeleteIcon")[0]);
        // delete second task
        fetch.mockResolvedValueOnce({ ok: true });
        fireEvent.click(screen.getAllByTestId("DeleteIcon")[1]);

        await waitFor(() => expect(screen.queryByText("Task 1")).not.toBeInTheDocument());
        await waitFor(() => expect(screen.queryByText("Task 2")).not.toBeInTheDocument());
        expect(screen.getByText(noTasksMessage)).toBeInTheDocument();
    });

    test("calls handleTaskCreated and adds a new task to the list", async () => {
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => [baseTask],
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));

        fireEvent.click(screen.getByTestId("AddCircleIcon"));

        fetch.mockResolvedValueOnce({ ok: true, json: async () => mockNewTask });

        fireEvent.change(screen.getByLabelText(/Title/i), { target: { value: mockNewTask.title } });
        fireEvent.change(screen.getByLabelText(/Description/i), { target: { value: mockNewTask.description } });
        fireEvent.change(screen.getByLabelText(/Status/i), { target: { value: mockNewTask.status } });
        fireEvent.click(screen.getByText("Submit"));

        await waitFor(() => screen.getByText("New Task"));

        expect(screen.getByText("New Task")).toBeInTheDocument();
        expect(screen.getByText("New description")).toBeInTheDocument();
        expect(screen.getByText("In Progress")).toBeInTheDocument();
    });

    test("calls handleTaskUpdated and updates the task in the list", async () => {
        const task = baseTask
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => [task],
        });

        render(<TasksTable />);
        await waitFor(() => screen.getByText("Task 1"));

        fireEvent.click(screen.getByTestId("EditIcon"));

        const updatedTask = { id: 1, title: "Updated Task", description: "Updated Description", status: "Completed" };
        fetch.mockResolvedValueOnce({ ok: true, json: async () => updatedTask });

        fireEvent.change(screen.getByLabelText(/Title/i), { target: { value: updatedTask.title } });
        fireEvent.change(screen.getByLabelText(/Description/i), { target: { value: updatedTask.description } });
        fireEvent.change(screen.getByLabelText(/Status/i), { target: { value: updatedTask.status } });
        fireEvent.click(screen.getByText("Submit"));

        await waitFor(() => screen.getByText(updatedTask.title));
        expect(screen.getByText(updatedTask.title)).toBeInTheDocument();
        expect(screen.getByText(updatedTask.description)).toBeInTheDocument();
        expect(screen.getByText(updatedTask.status)).toBeInTheDocument();

        expect(screen.queryByText(task.title)).not.toBeInTheDocument();
        expect(screen.queryByText(task.description)).not.toBeInTheDocument();
        expect(screen.queryByText(task.status)).not.toBeInTheDocument();
    });
});
