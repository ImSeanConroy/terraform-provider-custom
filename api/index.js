const express = require("express");
const { v4: uuidv4 } = require("uuid");

const app = express();
const PORT = 3000;

app.use(express.json()); // Middleware to parse JSON request bodies

// In-memory storage for notes
let notes = [
  {
    id: "3f5b4d2e-8c12-4b8b-8cbb-3df2bcad50e5",
    text: "Learn Robotics and SLAM",
    createdAt: "2025-02-03T10:00:00.000Z",
    updatedAt: "2025-02-03T10:05:00.000Z",
  },
];

// Route to get all notes
app.get("/notes", (req, res) => {
  res.json({ notes });
});

// Route to add a new note
app.post("/notes", (req, res) => {
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
});

// Route to update a note by ID
app.put("/notes/:id", (req, res) => {
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
});

// Route to delete a note by ID
app.delete("/notes/:id", (req, res) => {
  const { id } = req.params;
  const noteIndex = notes.findIndex((note) => note.id === id);

  if (noteIndex === -1) {
    return res.status(404).json({ error: "Note not found" });
  }

  const deletedNote = notes.splice(noteIndex, 1);
  res.json({ message: "Note deleted", note: deletedNote[0] });
});

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}`);
});
