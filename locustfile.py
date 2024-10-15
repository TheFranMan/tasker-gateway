from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def hello_world(self):
        self.client.get("/status/7fbef510-e37d-4884-97e2-c31fac6a89ae", headers={"Authorization": "auth-token-valid-1"})