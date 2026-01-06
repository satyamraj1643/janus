
import psycopg2
import sys

conn_str = "postgresql://janus_admin:t6exln67eDs6ygXQQTkaJE2CSsBvtMtl@dpg-d571ksuuk2gs73cp2sjg-a.oregon-postgres.render.com/janus_db_03vr"

try:
    conn = psycopg2.connect(conn_str)
    cur = conn.cursor()

    tables = ['jobs', 'batch']

    with open('schema_output.txt', 'w', encoding='utf-8') as f:
        # Get all tables
        cur.execute("SELECT table_name FROM information_schema.tables WHERE table_schema='public';")
        tables = [row[0] for row in cur.fetchall()]
        
        for table in tables:
            f.write(f"--- Schema for table: {table} ---\n")
            cur.execute(f"""
                SELECT column_name, data_type, is_nullable
                FROM information_schema.columns
                WHERE table_name = '{table}'
                ORDER BY ordinal_position;
            """)
            rows = cur.fetchall()
            for row in rows:
                f.write(f"{row[0]} ({row[1]}) - Nullable: {row[2]}\n")
            f.write("\n")
    
    cur.close()
    conn.close()

except Exception as e:
    print(f"Error: {e}")
