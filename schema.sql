
-- Users table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE
);

-- Songs table (unique by Spotify ID if you want to avoid duplicates)
CREATE TABLE songs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    artist TEXT,
    spotify_id TEXT UNIQUE
);

-- Each list (like "Top 50 – August 2025") belongs to one user
CREATE TABLE user_song_lists (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    list_name TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Each entry in a list = one song in a specific position
CREATE TABLE user_song_list_entries (
    list_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    position INTEGER NOT NULL, -- 1 through 50

    PRIMARY KEY (list_id, song_id), -- prevents duplicates in same list
    FOREIGN KEY (list_id) REFERENCES user_song_lists(id) ON DELETE CASCADE,
    FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

-- Helpful index for fast lookups by user
CREATE INDEX idx_user_lists ON user_song_lists(user_id);

-- Helpful index for fast lookups by song across lists
CREATE INDEX idx_song_entries ON user_song_list_entries(song_id);
