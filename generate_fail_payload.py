
import json
import uuid

# Configuration
TENANT_ID = "2dad64a8-3f87-4d6e-9b4c-1cfa5917fd4b" # Valid admin user
FAIL_PRIORITY = 0 # Assuming MinPriority > 0 in default config

def generate_fail_batch():
    batch_name = "failing_batch"
    jobs = []
    
    # Generate 100 jobs with low priority
    for i in range(100):
        jobs.append({
            "job_id": f"fail_job_{uuid.uuid4()}",
            "tenant_id": TENANT_ID,
            "priority": FAIL_PRIORITY,
            "payload": {
                "action": "fail_test"
            },
            "dependencies": {
                "heavy_resource": 1000 # High cost to trigger cost limit if priority doesn't
            }
        })
        
    payload = {
        "batch_name": batch_name,
        "jobs": jobs
    }
    
    print(json.dumps(payload, indent=2))

if __name__ == "__main__":
    generate_fail_batch()
