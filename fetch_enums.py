
import psycopg2

conn_str = "postgresql://janus_admin:t6exln67eDs6ygXQQTkaJE2CSsBvtMtl@dpg-d571ksuuk2gs73cp2sjg-a.oregon-postgres.render.com/janus_db_03vr"

try:
    conn = psycopg2.connect(conn_str)
    cur = conn.cursor()

    cur.execute("SELECT t.typname, e.enumlabel FROM pg_type t JOIN pg_enum e ON t.oid = e.enumtypid ORDER BY t.typname, e.enumsortorder;")
    rows = cur.fetchall()
    
    current_type = None
    for row in rows:
        if row[0] != current_type:
            print(f"\nType: {row[0]}")
            current_type = row[0]
        print(f"- {row[1]}")

    cur.close()
    conn.close()

except Exception as e:
    print(f"Error: {e}")
