CREATE TABLE segments(
    seg_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    tag TEXT NOT NULL UNIQUE
);
CREATE TABLE users_segments(
    blg_id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL,
    seg_id INT NOT NULL,
    create_time TIMESTAMPTZ DEFAULT NOW(),
    remove_time TIMESTAMPTZ,
    FOREIGN KEY (seg_id) REFERENCES segments(seg_id)
);
