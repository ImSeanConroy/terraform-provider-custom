/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const { Router } = require("express");
const {
  getTask,
  getTasks,
  createTask,
  updateTask,
  deleteTask,
} = require("../controller/tasks.controller");

const router = Router();

router
  .get("/", getTasks)
  .get("/:id", getTask)
  .post("/", createTask)
  .put("/:id", updateTask)
  .delete("/:id", deleteTask);

module.exports = router;
