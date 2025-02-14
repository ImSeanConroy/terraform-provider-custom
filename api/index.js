/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const express = require("express");
const tasksRouter = require("./routes/tasks.route");

const app = express();
const PORT = 3000;

app.use(express.json());

app.use("/tasks", tasksRouter);

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}/tasks`);
});
