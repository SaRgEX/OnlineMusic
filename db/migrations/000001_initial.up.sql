CREATE TABLE music_performer
(
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE song
(
    id SERIAL NOT NULL PRIMARY KEY,
    name TEXT NOT NULL,
    lyric TEXT,
    release_date DATE NOT NULL DEFAULT now(),
    link TEXT,
    performer_id INTEGER NOT NULL,
    FOREIGN KEY (performer_id) REFERENCES music_performer(id)
);

CREATE INDEX ix_song_performer_id ON song (performer_id);
CREATE INDEX ix_song_creation_date ON song (release_date);
CREATE INDEX ix_song_lyrics ON song USING gin(to_tsvector('english', lyric));