const { v4: uuidv4 } = require("uuid");
const notes = require("../data/notes")

// @description   Get all notes
// @route         POST /notes
// @access        Public
const getNotes = (req, res) => {
  res.json({ notes });
};

// @description   Create note
// @route         POST /notes
// @access        Public
const createNote = (req, res) => {
  const { text } = req.body;
  if (!text) {
    return res.status(400).json({ error: "Note text is required" });
  }

  const newNote = {
    id: uuidv4(),
    text,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  };

  notes.push(newNote);
  res.status(201).json({ message: "Note added", note: newNote });
};

// @description   Get note by ID
// @route         POST /notes/:id
// @access        Public
const getNote = (req, res) => {
  const { id } = req.params;
  const note = notes.find((note) => note.id === id);

  if (!note) {
    return res.status(404).json({ error: "Note not found" });
  }

  res.json(note);
};

// @description   Update note by ID
// @route         POST /notes/:id
// @access        Public
const updateNote = (req, res) => {
  const { id } = req.params;
  const { text } = req.body;
  const note = notes.find((note) => note.id === id);

  if (!note) {
    return res.status(404).json({ error: "Note not found" });
  }

  if (!text) {
    return res.status(400).json({ error: "Note text is required for update" });
  }

  note.text = text;
  note.updatedAt = new Date().toISOString();

  res.json({ message: "Note updated", note });
};

// @description   Delete note by ID
// @route         POST /notes/:id
// @access        Public
const deleteNote = (req, res) => {
  const { id } = req.params;
  const noteIndex = notes.findIndex((note) => note.id === id);

  if (noteIndex === -1) {
    return res.status(404).json({ error: "Note not found" });
  }

  const deletedNote = notes.splice(noteIndex, 1);
  res.json({ message: "Note deleted", note: deletedNote[0] });
};

module.exports = { getNotes, createNote, getNote, updateNote, deleteNote };
