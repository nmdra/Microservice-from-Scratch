import os
import time
import logging
from flask import Flask, request, jsonify
import psycopg2

app = Flask(__name__)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("payment-service")


def get_db():
    return psycopg2.connect(
        host=os.getenv("POSTGRES_HOST"),
        port=os.getenv("POSTGRES_PORT"),
        user=os.getenv("POSTGRES_USER"),
        password=os.getenv("POSTGRES_PASSWORD"),
        database=os.getenv("PAYMENT_DB"),
    )


# DB Retry
conn = None
for i in range(5):
    try:
        conn = get_db()
        break
    except Exception:
        logger.warning("Database not ready, retrying...")
        time.sleep(3)

if conn is None:
    logger.error("Could not connect to database")
    exit(1)

cur = conn.cursor()

cur.execute("""
CREATE TABLE IF NOT EXISTS payments(
 id SERIAL PRIMARY KEY,
 order_id INT,
 amount FLOAT,
 method TEXT,
 status TEXT
)
""")

conn.commit()

logger.info("Payment Service started")

# Health
@app.route("/health")
def health():
    try:
        cur.execute("SELECT 1")
        return {"status": "UP"}
    except:
        return {"status": "DOWN"}, 500


# Get Payments
@app.route("/payments")
def get_payments():
    cur.execute("SELECT * FROM payments")
    rows = cur.fetchall()

    result = []
    for r in rows:
        result.append(
            {
                "id": r[0],
                "orderId": r[1],
                "amount": r[2],
                "method": r[3],
                "status": r[4],
            }
        )

    return jsonify(result)


# Process Payment
@app.route("/payments/process", methods=["POST"])
def process_payment():
    data = request.json

    cur.execute(
        """
        INSERT INTO payments(order_id,amount,method,status)
        VALUES (%s,%s,%s,%s)
        RETURNING id
        """,
        (data["orderId"], data["amount"], data["method"], "SUCCESS"),
    )

    pid = cur.fetchone()[0]
    conn.commit()

    logger.info("Payment processed", extra={"payment_id": pid})

    return {"id": pid, "status": "SUCCESS"}, 201


# GET Payment
@app.route("/payments/<id>")
def get_payment(id):
    cur.execute("SELECT * FROM payments WHERE id=%s", (id,))
    row = cur.fetchone()

    if not row:
        return {"error": "Payment not found"}, 404

    return {
        "id": row[0],
        "orderId": row[1],
        "amount": row[2],
        "method": row[3],
        "status": row[4],
    }

if __name__ == "__main__":
    app.run(host="0.0.0.0", port=8080)
