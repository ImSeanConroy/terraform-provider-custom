/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const { v4: uuidv4 } = require("uuid");
const tasks = require("../data/tasks");

// @description   Get all tasks
// @route         POST /tasks
// @access        Public
const getTasks = (req, res) => {
  res.json({ tasks });
};

// @description   Create task
// @route         POST /tasks
// @access        Public
const createTask = (req, res) => {
  const { title, description, complete } = req.body;
  if (!title) {
    return res.status(400).json({ error: "task title is required" });
  }

  const newTask = {
    id: uuidv4(),
    title,
    description: description ? description : "",
    complete: complete ? complete : false,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  tasks.push(newTask);
  res.status(201).json({ message: "Task created", task: newTask });
};

// @description   Get task by ID
// @route         POST /tasks/:id
// @access        Public
const getTask = (req, res) => {
  const { id } = req.params;
  const task = tasks.find((task) => task.id === id);

  if (!task) {
    return res.status(404).json({ error: "Task not found" });
  }

  res.json(task);
};

// @description   Update task by ID
// @route         POST /tasks/:id
// @access        Public
const updateTask = (req, res) => {
  const { id } = req.params;
  const { title, description, complete } = req.body;
  const task = tasks.find((task) => task.id === id);

  if (!task) {
    return res.status(404).json({ error: "Task not found" });
  }

  if (!title) {
    return res.status(400).json({ error: "Task title is required for update" });
  }

  task.title = title;
  task.description = description ? description : "";
  task.complete = complete ? complete : false;
  task.updatedAt = new Date().toISOString();

  res.json({ message: "Task updated", task });
};

// @description   Delete task by ID
// @route         POST /tasks/:id
// @access        Public
const deleteTask = (req, res) => {
  const { id } = req.params;
  const taskIndex = tasks.findIndex((task) => task.id === id);

  if (taskIndex === -1) {
    return res.status(404).json({ error: "Task not found" });
  }

  const deletedTask = tasks.splice(taskIndex, 1);
  res.json({ message: "Task deleted", task: deletedTask[0] });
};

module.exports = { getTasks, createTask, getTask, updateTask, deleteTask };
