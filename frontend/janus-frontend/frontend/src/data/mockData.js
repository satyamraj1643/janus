// Simulating backend response for keys
export const mockStats = {
  total_requests: 12543,
  admitted_requests: 9876,
  rejected_requests: 2667,
  rejection_reasons: {
    "rate_limit_exceeded": 1543,
    "priority_too_low": 890,
    "duplicate_request": 234
  },
  recent_decisions: Array.from({ length: 50 }, (_, i) => ({
    job_id: `job_${1000 + i}`,
    admitted: Math.random() > 0.3,
    reason: Math.random() > 0.3 ? "admitted" : ["rate_limit_exceeded", "priority_too_low", "duplicate_request"][Math.floor(Math.random() * 3)],
    timestamp: new Date(Date.now() - i * 1000 * 60).toISOString(), // decreasing timestamps
    payload: {
      task_type: ["data_processing", "image_resize", "email_notification"][Math.floor(Math.random() * 3)],
      parameters: {
        priority: Math.floor(Math.random() * 100),
        region: ["us-east-1", "eu-west-1", "ap-south-1"][Math.floor(Math.random() * 3)],
        retries: Math.floor(Math.random() * 5)
      },
      metadata: {
        source: "api_gateway",
        attempt: 1
      }
    }
  }))
};
