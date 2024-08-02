CREATE TABLE IF NOT EXISTS script_recipes (
    script_id INTEGER REFERENCES scripts(id) ON DELETE CASCADE,
    recipe_id INTEGER REFERENCES recipes(id) ON DELETE CASCADE,
    PRIMARY KEY (script_id, recipe_id)
);