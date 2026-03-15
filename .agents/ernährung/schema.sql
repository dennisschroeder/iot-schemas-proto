CREATE TABLE IF NOT EXISTS daily_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL UNIQUE,
    day_number INTEGER,
    calories_target INTEGER,
    protein_target INTEGER,
    fat_target INTEGER,
    carbs_target INTEGER,
    calories_actual INTEGER,
    protein_actual INTEGER,
    fat_actual INTEGER,
    carbs_actual INTEGER,
    water_liters REAL,
    training_info TEXT,
    supplements TEXT,
    conclusion TEXT,
    status TEXT DEFAULT 'active' -- e.g., 'active', 'recovery'
);

CREATE TABLE IF NOT EXISTS weight_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date TEXT NOT NULL UNIQUE,
    weight REAL NOT NULL,
    comment TEXT
);

CREATE TABLE IF NOT EXISTS meals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    log_id INTEGER,
    meal_type TEXT, -- 'breakfast', 'lunch', 'dinner', 'snack', 'post-workout', 'late-night'
    description TEXT,
    calories INTEGER,
    protein INTEGER,
    fat INTEGER,
    carbs INTEGER,
    FOREIGN KEY(log_id) REFERENCES daily_logs(id)
);

CREATE TABLE IF NOT EXISTS foods (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    brand TEXT,
    calories_per_100 INTEGER,
    protein_per_100 REAL,
    fat_per_100 REAL,
    carbs_per_100 REAL,
    serving_size_g INTEGER,
    source TEXT, -- z.B. 'Verpackung', 'FDDB', 'ChatGPT Schätzung'
    last_used TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
