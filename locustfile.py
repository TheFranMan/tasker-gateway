from locust import HttpUser, task, constant_throughput

class StatusEndpointUser(HttpUser):
    wait_time = constant_throughput(100)
    fixed_count = 1000

    @task(10)
    def token_exists(self):
        self.client.get("/status/e96b72b8-fe24-46b8-8525-280fac1032fd", headers={"Authorization": "auth-token-valid-1"}, name="status")

    @task(1)
    def token_doeas_not_exist(self):
        self.client.get("/status/7fbef510-e37d-4884-97e2-c31fac6a89aa", headers={"Authorization": "auth-token-valid-1"}, name="status")