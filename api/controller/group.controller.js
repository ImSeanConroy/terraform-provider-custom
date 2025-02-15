/**
 * Copyright (c) HashiCorp, Inc.
 * SPDX-License-Identifier: MPL-2.0
 */

const { v4: uuidv4 } = require("uuid");
const groups = require("../data/groups");

// @description   Get all groups
// @route         POST /groups
// @access        Public
const getGroups = (req, res) => {
  res.json({ groups });
};

// @description   Create group
// @route         POST /groups
// @access        Public
const createGroup = (req, res) => {
  const { title } = req.body;
  if (!title) {
    return res.status(400).json({ error: "Group title is required" });
  }

  const newGroup = {
    id: uuidv4(),
    title,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  groups.push(newGroup);
  res.status(201).json({ message: "Group created", group: newGroup });
};

// @description   Get group by ID
// @route         POST /groups/:id
// @access        Public
const getGroup = (req, res) => {
  const { id } = req.params;
  const group = groups.find((group) => group.id === id);

  if (!group) {
    return res.status(404).json({ error: "Group not found" });
  }

  res.json(group);
};

// @description   Update group by ID
// @route         POST /groups/:id
// @access        Public
const updateGroup = (req, res) => {
  const { id } = req.params;
  const { title } = req.body;
  const group = groups.find((group) => group.id === id);

  if (!group) {
    return res.status(404).json({ error: "Group not found" });
  }

  if (!title) {
    return res.status(400).json({ error: "Group title is required for update" });
  }

  group.title = title;
  group.updatedAt = new Date().toISOString();

  res.json({ message: "Group updated", group });
};

// @description   Delete group by ID
// @route         POST /groups/:id
// @access        Public
const deleteGroup = (req, res) => {
  const { id } = req.params;
  const groupindex = groups.findIndex((group) => group.id === id);

  if (groupindex === -1) {
    return res.status(404).json({ error: "Group not found" });
  }

  const deletedGroup = groups.splice(groupindex, 1);
  res.json({ message: "Group deleted", group: deletedGroup[0] });
};

module.exports = { getGroups, createGroup, getGroup, updateGroup, deleteGroup };
