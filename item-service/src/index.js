import express from "express";
import pkg from "pg";

const { Pool } = pkg;

const app = express();
app.use(express.json());

const pool = new Pool({
    host: process.env.POSTGRES_HOST,
    port: process.env.POSTGRES_PORT,
    user: process.env.POSTGRES_USER,
    password: process.env.POSTGRES_PASSWORD,
    database: process.env.ITEM_DB,
});

async function initDB() {
    await pool.query(`
    CREATE TABLE IF NOT EXISTS items(
      id SERIAL PRIMARY KEY,
      name TEXT NOT NULL
    )
  `);
}

await initDB();

app.get("/items", async (req, res) => {
    const result = await pool.query("SELECT * FROM items");
    res.json(result.rows);
})

app.post("/items", async (req, res) => {
    const { name } = req.body;

    const result = await pool.query(
        "INSERT INTO items(name) VALUES($1) RETURNING *",
        [name]
    );

    res.status(201).json(result.rows[0]);
});

app.get("/items/:id", async (req, res) => {
    const result = await pool.query(
        "SELECT * FROM items WHERE id=$1",
        [req.params.id]
    );

    res.json(result.rows[0]);
});

app.listen(8080, () => {
    console.log("Item Service running on 8080");
});