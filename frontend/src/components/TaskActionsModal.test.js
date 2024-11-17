import React from "react";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import TaskActionsModal from "./TaskActionsModal";

jest.mock("node-fetch", () => jest.fn());
global.fetch = jest.fn();
global.alert = jest.fn()

describe("TaskActionsModal", () => {
    const mockOnClose = jest.fn();
    const mockOnTaskUpdated = jest.fn();
    const mockOnTaskCreated = jest.fn();
    const labelTitle = /Title/i
    const labelDescription = /Description/i
    const labelStatus = /Status/i
    const baseTask = {
        id: 1,
        title: "Test Task",
        description: "Test Description",
        status: "Pending",
    };

    const updatedTask = {
        id: 1,
        title: "Updated Task",
        description: "Updated Description",
        status: "Done",
    };

    beforeEach(() => {
        jest.clearAllMocks();
    });

    it("does nothing when open is false", () => {
        render(
            <TaskActionsModal
                open={false}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        expect(screen.queryByText("Create Task")).not.toBeInTheDocument();
        expect(screen.queryByLabelText("Title")).not.toBeInTheDocument();
        expect(screen.queryByLabelText("Description")).not.toBeInTheDocument();
        expect(screen.queryByLabelText("Status")).not.toBeInTheDocument();
        expect(screen.queryByText("Cancel")).not.toBeInTheDocument();
        expect(screen.queryByText("Submit")).not.toBeInTheDocument();

        expect(mockOnClose).not.toHaveBeenCalled();
        expect(mockOnTaskUpdated).not.toHaveBeenCalled();
        expect(mockOnTaskCreated).not.toHaveBeenCalled();
    });

    it("renders correctly in create mode", () => {
        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        expect(screen.getByText("Create Task")).toBeInTheDocument();
        expect(screen.getByLabelText(labelTitle)).toBeInTheDocument();
        expect(screen.getByLabelText(labelDescription)).toBeInTheDocument();
        expect(screen.getByLabelText(labelStatus)).toBeInTheDocument();
        expect(screen.getByText("Cancel")).toBeInTheDocument();
        expect(screen.getByText("Submit")).toBeInTheDocument();
    });

    it("renders correctly in update mode - defaults are in", () => {
        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={baseTask}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        expect(screen.getByText(`Edit ${baseTask.title}`)).toBeInTheDocument();
        expect(screen.getByDisplayValue(baseTask.title)).toBeInTheDocument();
        expect(screen.getByDisplayValue(baseTask.description)).toBeInTheDocument();
        expect(screen.getByDisplayValue(baseTask.status)).toBeInTheDocument();
    });

    it("validates empty title on submission", async () => {
        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        fireEvent.click(screen.getByText("Submit"));
        await waitFor(() => {
            expect(screen.getByText("Title cannot be empty")).toBeInTheDocument();
        });
    });

    it("succeeds to submit new task and reset form labels", async () => {
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => (baseTask),
        });

        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        fireEvent.change(screen.getByLabelText(labelTitle), { target: { value: "New Task" } });
        fireEvent.change(screen.getByLabelText(labelDescription), { target: { value: "Task Description" } });
        fireEvent.change(screen.getByLabelText(labelStatus), { target: { value: "Pending" } });

        fireEvent.click(screen.getByText("Submit"));

        await waitFor(() => {
            expect(fetch).toHaveBeenCalledWith("http://localhost:8080/tasks", expect.anything());
            expect(mockOnTaskCreated).toHaveBeenCalledWith({
                id: 1,
                title: "Test Task",
                description: "Test Description",
                status: "Pending",
            });
            expect(mockOnClose).toHaveBeenCalled();
        });
    });

    it("updates task correctly", async () => {
        const oldTask = { id: 1, title: "Old Task", description: "Old Description", status: "In Progress" };
        fetch.mockResolvedValueOnce({
            ok: true,
            json: async () => (updatedTask),
        });

        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={oldTask}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        fireEvent.change(screen.getByLabelText(labelTitle), { target: { value: "Updated Task" } });
        fireEvent.change(screen.getByLabelText(labelDescription), { target: { value: "Updated Description" } });
        fireEvent.change(screen.getByLabelText(labelStatus), { target: { value: "Done" } });

        fireEvent.click(screen.getByText("Submit"));

        await waitFor(() => {
            expect(fetch).toHaveBeenCalledWith(
                `http://localhost:8080/tasks/${oldTask.id}`,
                expect.anything()
            );
            expect(mockOnTaskUpdated).toHaveBeenCalledWith(updatedTask);
            expect(mockOnClose).toHaveBeenCalled();
        });
    });

    it("popping an alert when submission error", async () => {
        fetch.mockResolvedValueOnce({
            ok: false,
            statusText: "Internal Server Error",
        });

        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        fireEvent.change(screen.getByLabelText(labelTitle), { target: { value: "New Task" } });
        fireEvent.change(screen.getByLabelText(labelDescription), { target: { value: "Task Description" } });
        fireEvent.change(screen.getByLabelText(labelStatus), { target: { value: "Pending" } });

        fireEvent.click(screen.getByText("Submit"));

        await waitFor(() => expect(global.alert).toHaveBeenCalledWith('Error: Internal Server Error'));

        expect(global.alert).toHaveBeenCalledTimes(1);
        expect(global.alert).toHaveBeenCalledWith('Error: Internal Server Error');

    });

    it("calls onClose when cancel is clicked", () => {
        render(
            <TaskActionsModal
                open={true}
                onClose={mockOnClose}
                task={null}
                onTaskUpdated={mockOnTaskUpdated}
                onTaskCreated={mockOnTaskCreated}
            />
        );

        fireEvent.click(screen.getByText("Cancel"));

        expect(mockOnClose).toHaveBeenCalled();
    });
});
