const { Router } = require("express");
const {
  getNotes,
  getNote,
  createNote,
  updateNote,
  deleteNote,
} = require("../controller/notes.controller");

const router = Router();

router
  .get("/", getNotes)
  .get("/:id", getNote)
  .post("/", createNote)
  .put("/:id", updateNote)
  .delete("/:id", deleteNote);

module.exports = router;
