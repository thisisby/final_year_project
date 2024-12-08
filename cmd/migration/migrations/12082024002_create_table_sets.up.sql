CREATE TABLE IF NOT EXISTS exercise_sets (
    id SERIAL PRIMARY KEY,
    reps INT NOT NULL,
    weight DECIMAL(5, 2) NOT NULL,
    notes TEXT DEFAULT '',
    workout_exercise_id INT NOT NULL REFERENCES workout_exercises(id) ON DELETE CASCADE,
    owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL
);


