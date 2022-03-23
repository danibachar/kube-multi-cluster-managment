from locust import HttpUser, task

class HelloWorldUser(HttpUser):
    @task
    def hello_world(self):
        self.client.post("/load", json={"memory_params": {"duration_seconds": 0.2, "kb_count": 50}, "cpu_params": {"duration_seconds": 0.2, "load": 0.2}})