
import psycopg2

conn_str = "postgresql://janus_admin:t6exln67eDs6ygXQQTkaJE2CSsBvtMtl@dpg-d571ksuuk2gs73cp2sjg-a.oregon-postgres.render.com/janus_db_03vr"

try:
    conn = psycopg2.connect(conn_str)
    cur = conn.cursor()

    print("Fetching users...")
    cur.execute("SELECT user_id, name FROM users;")
    rows = cur.fetchall()
    for row in rows:
        print(f"User: {row[0]} ({row[1]})")
    
    cur.close()
    conn.close()

except Exception as e:
    print(f"Error: {e}")
