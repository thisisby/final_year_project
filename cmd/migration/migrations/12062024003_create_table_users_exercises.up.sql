CREATE TABLE IF NOT EXISTS workout_exercises (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL DEFAULT '',
    main_note TEXT NOT NULL DEFAULT '',
    secondary_note TEXT NOT NULL DEFAULT '',
    workout_id INT NOT NULL REFERENCES workouts(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS idx_workout_exercises_name ON workout_exercises(name);


