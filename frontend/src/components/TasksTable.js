import * as React from 'react';
import { styled } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell, { tableCellClasses } from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import Paper from '@mui/material/Paper';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import AddCircleIcon from '@mui/icons-material/AddCircle';

import './TasksTable.css'
import {useEffect, useState} from "react";
import {CircularProgress} from "@mui/material";
import TaskActionsModal from "./TaskActionsModal";

const StyledTableCell = styled(TableCell)(({ theme }) => ({
    [`&.${tableCellClasses.head}`]: {
        backgroundColor: theme.palette.common.black,
        color: theme.palette.common.white,
    },
    [`&.${tableCellClasses.body}`]: {
        fontSize: 14,
    },
}));

const StyledTableRow = styled(TableRow)(({ theme }) => ({
    '&:nth-of-type(odd)': {
        backgroundColor: theme.palette.action.hover,
    },
    // hide last border
    '&:last-child td, &:last-child th': {
        border: 0,
    },
}));

export default function TasksTable() {
    const [tasks, setTasks] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedTask, setSelectedTask] = useState(null);
    const [modalOpen, setModalOpen] = useState(false);


    useEffect(() => {
        const fetchTasks = async () => {
            try {
                setLoading(true); // Set loading to while fetching data
                const response = await fetch("http://localhost:8080/tasks");

                if (!response.ok) {
                    throw new Error(`HTTP error! Status: ${response.status}`);
                }

                const data = await response.json();
                setTasks(data);
            } catch (err) {
                alert(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchTasks();
    }, []);

    const editTask = (task) => {
        setSelectedTask(task);
        setModalOpen(true);
    }

    const handleCloseModal = () => {
        setModalOpen(false);
        setSelectedTask(null);
    };

    const handleTaskUpdated = (updatedTask) => {
        setTasks((prevTasks) =>
            prevTasks.map((task) =>
                task.id === updatedTask.id ? { ...task, ...updatedTask } : task
            )
        );
    };

    const handleTaskCreated = (newTask) => {
        setTasks((prevTasks) => [...prevTasks, newTask]);
    };

    const addNewTask = () => {
        setModalOpen(true);
        setSelectedTask(null);
    }

    const deleteTask = async (taskId) => {
        try {
            const response = await fetch(`http://localhost:8080/tasks/${taskId}`, {
                method: 'DELETE',
                headers: {
                    'Content-Type': 'application/json',
                },
                credentials: 'include',
            });

            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            setTasks((prevTasks) => prevTasks.filter((task) => task.id !== taskId));
        } catch (err) {
            alert("Failed to delete task:", err.message)
        }
    };

    console.log(tasks)

    return (
        <>
            {loading ? <CircularProgress />:  (
                <TableContainer className="tasks-table" component={Paper}>
                    <Table sx={{ minWidth: 700 }} aria-label="customized table">
                        <TableHead>
                            <TableRow>
                                <StyledTableCell>Title</StyledTableCell>
                                <StyledTableCell align="right">Description</StyledTableCell>
                                <StyledTableCell align="right">Status</StyledTableCell>
                                <StyledTableCell align="right"><AddCircleIcon fontSize="large" color="success" onClick={() => addNewTask()}/></StyledTableCell>
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {tasks?.map((task) => (
                                <StyledTableRow key={task.title}>
                                    <StyledTableCell component="th" scope="row">
                                        {task.title}
                                    </StyledTableCell>
                                    <StyledTableCell align="right">{task.description}</StyledTableCell>
                                    <StyledTableCell align="right">{task.status}</StyledTableCell>
                                    <StyledTableCell align="left">
                                        <div className="tasks-actions">
                                            <EditIcon onClick={() => editTask(task)} />
                                            <DeleteIcon onClick={() => deleteTask(task.id)} />
                                        </div>
                                    </StyledTableCell>
                                </StyledTableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            )}
            {tasks?.length === 0 ? <span className="no-tasks-label">No tasks to show, use the create (+) button to add tasks.</span> : <></>}
            <TaskActionsModal open={modalOpen} onClose={handleCloseModal} task={selectedTask} onTaskUpdated={handleTaskUpdated} onTaskCreated={handleTaskCreated} />
        </>

    );
}