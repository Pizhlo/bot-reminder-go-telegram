ALTER TABLE notes.notes ADD
last_edit timestamp;

INSERT INTO users.state_types VALUES
(24, 'edit_note');