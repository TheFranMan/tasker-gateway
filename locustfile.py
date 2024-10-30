from locust import HttpUser, task, events
from itertools import count
from typing import Final
import requests

DOMAIN: Final[str] = "http://localhost:3005"
PATH_API_POLL: Final[str] = "/api/poll"
PATH_API_DELETE: Final[str] = "/api/user"
HEADERS: Final[dict] = {"Authorization": "auth-token-valid-1", "Content-Type": "application/json"}

test_token: str = ""

@events.test_start.add_listener
def on_test_start(environment, **kwargs) -> None:
    response = requests.post(DOMAIN + PATH_API_DELETE, headers=HEADERS, json = {"id": 1})
    global test_token
    test_token = response.json()["token"]

class PollEndpoint(HttpUser):
    @task(10)
    def token_exists(self):
        global test_token
        self.client.get(DOMAIN + PATH_API_POLL + "/" + test_token, headers=HEADERS, name="/api/poll")

    @task(1)
    def token_does_not_exist(self):
        self.client.get(DOMAIN + PATH_API_POLL + "/88888888-4444-4444-4444-121212121212", headers=HEADERS, name="/api/poll")

class DeleteEndpoint(HttpUser):
    _ids = count(1)

    def __init__(self, parent) -> None:
        super().__init__(parent)
        self.id = next(self._ids)

    @task()
    def delete_request(self) -> None:
        self.client.post(DOMAIN + PATH_API_DELETE, headers=HEADERS, json = {"id": self.id}, name="/api/user")