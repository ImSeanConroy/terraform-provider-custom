const express = require("express");
const notesRouter = require("./routes/notes.route");

const app = express();
const PORT = 3000;

app.use(express.json());

app.use("/notes", notesRouter);

// Start the server
app.listen(PORT, () => {
  console.log(`Server is running on http://localhost:${PORT}/notes`);
});
