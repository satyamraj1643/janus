
import json
import uuid

# Configuration
TENANT_ID = "2dad64a8-3f87-4d6e-9b4c-1cfa5917fd4b" # Valid admin user
SUCCESS_PRIORITY = 100 # High priority

def generate_success_batch():
    batch_name = "success_atomic_batch"
    jobs = []
    
    # Generate 5 jobs with high priority and low cost
    for i in range(5):
        jobs.append({
            "job_id": f"success_atomic_job_{uuid.uuid4()}",
            "tenant_id": TENANT_ID,
            "priority": SUCCESS_PRIORITY,
            "payload": {
                "action": "process_image",
                "image_url": f"http://example.com/image_{i}.jpg"
            },
            "dependencies": {
                "db_read": 1 # Low cost
            }
        })
        
    payload = {
        "batch_name": batch_name,
        "jobs": jobs
    }
    
    print(json.dumps(payload, indent=2))

if __name__ == "__main__":
    generate_success_batch()
