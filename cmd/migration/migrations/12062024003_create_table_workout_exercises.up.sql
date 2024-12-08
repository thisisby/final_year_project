CREATE TABLE IF NOT EXISTS workout_exercises (
    id SERIAL PRIMARY KEY,
    main_note TEXT NOT NULL DEFAULT '',
    secondary_note TEXT NOT NULL DEFAULT '',
    workout_id INT NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    owner_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id INT NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL,
    deleted_at TIMESTAMP DEFAULT NULL,

    UNIQUE (workout_id, exercise_id, owner_id)
);



