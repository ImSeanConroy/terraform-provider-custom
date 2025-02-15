/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const { Router } = require("express");
const {
  getGroup,
  getGroups,
  createGroup,
  updateGroup,
  deleteGroup,
} = require("../controller/group.controller");

const router = Router();

router
  .get("/", getGroups)
  .get("/:id", getGroup)
  .post("/", createGroup)
  .put("/:id", updateGroup)
  .delete("/:id", deleteGroup);

module.exports = router;
