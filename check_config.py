
import psycopg2
import json

conn_str = "postgresql://janus_admin:t6exln67eDs6ygXQQTkaJE2CSsBvtMtl@dpg-d571ksuuk2gs73cp2sjg-a.oregon-postgres.render.com/janus_db_03vr"

try:
    conn = psycopg2.connect(conn_str)
    cur = conn.cursor()

    print("Fetching global_job_config...")
    cur.execute("SELECT user_id, config FROM global_job_config;")
    rows = cur.fetchall()
    for row in rows:
        print(f"User ID from Config Table: {row[0]}")
        # print(f"Config JSON: {json.dumps(row[1], indent=2)}") 
    
    cur.close()
    conn.close()

except Exception as e:
    print(f"Error: {e}")
