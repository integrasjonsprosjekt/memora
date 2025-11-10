import os
from locust import HttpUser, task, between, events
import json
import random
import string
from datetime import datetime

def load_tokens():
    if os.path.exists("stress_test_users.json"):
        with open("stress_test_users.json", "r") as f:
            users = json.load(f)
            tokens = [user['token'] for user in users]
            return tokens
    return []

TOKENS = load_tokens()

class MemoraUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        self.user_id = None
        self.deck_ids = []
        self.card_ids = []
        self.auth_token = random.choice(TOKENS) if TOKENS else None
        
        self.create_test_user()
        
        self.create_decks_on_start(count=2)
    
    def create_test_user(self):
        payload = {
            "name": f"Test User {self.generate_random_string(5)}"
        }
        
        with self.client.post(
            "/api/v1/users/",
            json=payload,
            catch_response=True,
            name="Setup: Create User"
        ) as response:
            if response.status_code == 201:
                data = response.json()
                self.user_id = data.get("id")
                response.success()
            else:
                response.failure(f"Failed to create user: {response.text}")
    
    def create_decks_on_start(self, count=2):
        for _ in range(count):
            payload = {
                "name": f"Initial Deck {self.generate_random_string(5)}",
            }
            headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
            with self.client.post(
                "/api/v1/decks/",
                json=payload,
                headers=headers,
                catch_response=True,
                name="Setup: Create Initial Deck"
            ) as response:
                if response.status_code == 201:
                    data = response.json()
                    deck_id = data.get("id")
                    if deck_id:
                        self.deck_ids.append(deck_id)
                    response.success()
                else:
                    response.failure(f"Failed to create initial deck: {response.text}")
    
    @staticmethod
    def generate_random_string(length=10):
        return ''.join(random.choices(string.ascii_letters + string.digits, k=length))
    
    @task(10)
    def health_check(self):
        with self.client.get("/status/", name="Health Check", catch_response=True) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Health check failed: {response.status_code}")
    
    @task(5)
    def get_user(self):
        if not self.user_id:
            return
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.get(
            f"/api/v1/users/",
            headers=headers,
            catch_response=True,
            name="GET User"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Get user failed: {response.status_code}")
    
    @task(5)
    def get_user_decks(self):
        if not self.user_id:
            return
    
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.get(
            f"/api/v1/users/decks/",
            headers=headers,
            catch_response=True,
            name="GET User Decks"
        ) as response:
            if response.status_code in [200, 404]:
                self.deck_ids = [deck['id'] for deck in response.json()] if response.status_code == 200 else []
                response.success()
            else:
                response.failure(f"Get user decks failed: {response.status_code}")
    
    @task(6)
    def create_deck(self):
        """Create a new deck."""
        payload = {
            "name": f"Deck {self.generate_random_string(5)}",
            "description": "A test deck created during stress testing"
        }
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.post(
            "/api/v1/decks/",
            json=payload,
            headers=headers,
            catch_response=True,
            name="POST Create Deck"
        ) as response:
            if response.status_code == 201:
                data = response.json()
                deck_id = data.get("id")
                if deck_id:
                    self.deck_ids.append(deck_id)
                response.success()
            else:
                response.failure(f"Create deck failed: {response.text}")  
                  
    @task(2)
    def update_deck(self):
        if not self.user_id:
            return
        
        deck_id = random.choice(self.deck_ids)
        payload = {
            "name": f"Updated Deck {self.generate_random_string(5)}"
        }
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        
        with self.client.get(
            f"/api/v1/decks/{deck_id}/",
            headers=headers,
            catch_response=True,
            name="GET Deck"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Get deck failed: {response.status_code}")
    
    @task(4)
    def create_card_in_deck(self):
        """Create a card in a deck."""
        if not self.deck_ids:
            return
        
        deck_id = random.choice(self.deck_ids)
        payload = {
            "type":"front_back",
            "front":"What is the capital of France?",
            "back":"Paris"
        }
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.post(
            f"/api/v1/decks/{deck_id}/cards/",
            json=payload,
            headers=headers,
            catch_response=True,
            name="POST Create Card in Deck"
        ) as response:
            if response.status_code == 201:
                data = response.json()
                card_id = data.get("id")
                if card_id:
                    self.card_ids.append(card_id)
                response.success()
            else:
                response.failure(f"Create card failed: {response.text}")
    
    @task(6)
    def get_card(self):
        if not self.card_ids or not self.deck_ids:
            return
        
        card_id = random.choice(self.card_ids)
        deck_id = random.choice(self.deck_ids)
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.get(
            f"/api/v1/decks/{deck_id}/cards/{card_id}/",
            headers=headers,
            catch_response=True,
            name="GET Card"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Get card failed: {response.status_code}")
    
    @task(2)
    def update_card(self):
        if not self.card_ids or not self.deck_ids:
            return
        
        card_id = random.choice(self.card_ids)
        deck_id = random.choice(self.deck_ids)
        payload = {
            "front": f"Updated Front {self.generate_random_string(5)}",
            "back": f"Updated Back {self.generate_random_string(5)}",
            "type": "front_back"
        }
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.put(
            f"/api/v1/decks/{deck_id}/cards/{card_id}/",
            json=payload,
            headers=headers,
            catch_response=True,
            name="PUT Update Card"
        ) as response:
            if response.status_code == 200:
                response.success()
            else:
                response.failure(f"Update card failed: {response.status_code}")
    @task(1)
    def delete_card(self):
        if not self.card_ids or not self.deck_ids:
            return
        
        card_id = random.choice(self.card_ids)
        deck_id = random.choice(self.deck_ids)
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.delete(
            f"/api/v1/decks/{deck_id}/cards/{card_id}/",
            headers=headers,
            catch_response=True,
            name="DELETE Card"
        ) as response:
            if response.status_code == 204:
                self.card_ids.remove(card_id)
                response.success()
            else:
                response.failure(f"Delete card failed: {response.status_code}")
    
@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print("\n" + "="*60)
    print("Starting Memora API Stress Test")
    print("="*60)
    print(f"Target: {environment.host}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("="*60 + "\n")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("\n" + "="*60)
    print("Stopping Memora API Stress Test")
    print("="*60)
    print(f"Target: {environment.host}")
    print(f"Time: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print("="*60 + "\n")
    
class ReadOnlyUser(HttpUser):
    wait_time = between(0.5, 2)
    weight = 3
    
    @task
    def read_health(self):
        self.client.get("/api/v1/status/", name="Health Check ReadOnly")

class WriteHeavyUser(HttpUser):
    """
    A user that primarily creates data.
    Useful for testing write performance.
    """
    wait_time = between(2, 5)
    weight = 1
    
    @task
    def create_many_decks(self):
        """Rapidly create multiple decks."""
        for _ in range(3):
            payload = {
                "name": f"Bulk Deck {MemoraUser.generate_random_string(5)}",
            }
            self.client.post("/api/v1/decks/", json=payload, name="WriteHeavy: Create Deck")