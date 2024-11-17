import React, { useState, useEffect } from "react";
import { Dialog, DialogActions, DialogContent, DialogTitle, TextField, Button, CircularProgress } from "@mui/material";

export default function TaskActionsModal({ open, onClose, task, onTaskUpdated, onTaskCreated }) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [status, setStatus] = useState("");
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState(null);

    // Set initial values in case of update
    useEffect(() => {
        if (task) {
            setTitle(task.title);
            setDescription(task.description);
            setStatus(task.status);
        }
    }, [task]);

    const isUpdateMode = task !== null;
    const handleCreation = () => {
        const task = {
            title,
            description,
            status,
        };

        return fetch(`http://localhost:8080/tasks`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(task),
        });
    }

    const handleUpdate = () => {
        const updatedTask = {
            title,
            description,
            status,
        };

        return fetch(`http://localhost:8080/tasks/${task.id}`, {
            method: "PUT",
                headers: {
            "Content-Type": "application/json",
        },
            body: JSON.stringify(updatedTask),
        });
    }

    const modalOnClose = () => {
        setTitle("");
        setDescription("");
        setStatus("");
        setError(null);

        onClose();
    }

    const handleSubmit = async () => {
        if (loading) return;

        if (!title || title.length === 0) {
            setError("Title cannot be empty");
            return;
        }

        try {
            setLoading(true);
            setError(null);

            let response;

            if (isUpdateMode) {
                response = await handleUpdate();
            } else {
                response = await handleCreation();
            }

            if (!response.ok) {
                throw new Error(`Error: ${response.statusText}`);
            }

            const updatedTaskData = await response.json(); // Get the new/updated task data from the server

            if (isUpdateMode) {
                onTaskUpdated(updatedTaskData);
            } else {
                onTaskCreated(updatedTaskData);
            }

            modalOnClose();
        } catch (error) {
            alert(error.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>{task ? `Edit ${task.title}` : "Create Task"}</DialogTitle>
            <DialogContent>
                <TextField
                    label="Title"
                    variant="outlined"
                    fullWidth
                    value={title}
                    required="true"
                    onChange={(e) => setTitle(e.target.value)}
                    margin="normal"
                />
                <TextField
                    label="Description"
                    variant="outlined"
                    fullWidth
                    value={description}
                    onChange={(e) => setDescription(e.target.value)}
                    margin="normal"
                />
                <TextField
                    label="Status"
                    variant="outlined"
                    fullWidth
                    value={status}
                    onChange={(e) => setStatus(e.target.value)}
                    margin="normal"
                />
                {error && <p style={{ color: "red" }}>{error}</p>}
            </DialogContent>
            <DialogActions>
                <Button onClick={modalOnClose} color="secondary">
                    Cancel
                </Button>
                <Button onClick={handleSubmit} color="primary" disabled={loading}>
                    {loading ? <CircularProgress size={24} /> : "Submit"}
                </Button>
            </DialogActions>
        </Dialog>
    );
}
