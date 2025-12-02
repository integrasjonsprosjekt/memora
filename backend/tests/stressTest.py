import os
import sys
from locust import HttpUser, task, between, events
import json
import random
import string
from datetime import datetime

def load_tokens():
    if os.path.exists("stress_test_users.json"):
        with open("stress_test_users.json", "r") as f:
            users = json.load(f)
            return users
    return []

def generate_random_string(length=10):
    return ''.join(random.choices(string.ascii_letters + string.digits, k=length))

USERS = load_tokens()

class MemoraUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        self.user_data = random.choice(USERS) if USERS else None
        if not USERS:
            print("No users available for authentication. Please run token generation script.")
            sys.exit(1)
        self.user_id = self.user_data.get("uid")
        self.deck_ids = []
        self.card_ids = []
        self.auth_token = self.user_data.get("token")
        self.decks = {}

        self.create_decks_on_start(count=2)
    
    def create_decks_on_start(self, count=2):
        if self.user_id is None:
            return
        for _ in range(count):
            payload = {
                "title": f"Initial Deck {generate_random_string(5)}",
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
                        self.decks[deck_id] = []
                    response.success()
                elif response.status_code == 429:
                    response.success()  # Rate limited, but don't fail
                else:
                    response.failure(f"Failed to create initial deck: {response.text}")
    
    @task(10)
    def health_check(self):
        with self.client.get("/api/v1/status/", name="Health Check", catch_response=True) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
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
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
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
            if response.status_code == 200:
                data = response.json()
                # API returns {"owned_decks": [...], "shared_decks": [...]}
                owned = data.get("owned_decks") or []
                shared = data.get("shared_decks") or []
                self.deck_ids = [deck['id'] for deck in owned + shared]
                response.success()
            elif response.status_code == 404:
                self.deck_ids = []
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Get user decks failed: {response.status_code}")
    
    @task(6)
    def create_deck(self):
        """Create a new deck."""
        payload = {
            "title": f"Deck {generate_random_string(5)}",
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
                    self.decks[deck_id] = []
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Create deck failed: {response.text}")  
                  
    @task(2)
    def update_deck(self):
        if not self.user_id or not self.deck_ids:
            return
        
        deck_id = random.choice(self.deck_ids)
        payload = {
            "title": f"Updated Deck {generate_random_string(5)}"
        }
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        
        with self.client.patch(
            f"/api/v1/decks/{deck_id}/",
            json=payload,
            headers=headers,
            catch_response=True,
            name="PATCH Update Deck"
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Update deck failed: {response.status_code}")
    
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
                    if deck_id in self.decks:
                        self.decks[deck_id].append(card_id)
                    else:
                        self.decks[deck_id] = [card_id]
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Create card failed: {response.text}")
    
    @task(6)
    def get_card(self):
        decks_with_cards = {deck_id: cards for deck_id, cards in self.decks.items() if cards}
        if not decks_with_cards:
            return

        deck_id = random.choice(list(decks_with_cards.keys()))
        card_id = random.choice(decks_with_cards[deck_id])
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.get(
            f"/api/v1/decks/{deck_id}/cards/{card_id}/",
            headers=headers,
            catch_response=True,
            name="GET Card"
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Get card failed: {response.status_code}")
    
    @task(2)
    def update_card(self):
        decks_with_cards = {deck_id: cards for deck_id, cards in self.decks.items() if cards}
        if not decks_with_cards:
            return

        deck_id = random.choice(list(decks_with_cards.keys()))
        card_id = random.choice(decks_with_cards[deck_id])
        payload = {
            "front": f"Updated Front {generate_random_string(5)}",
            "back": f"Updated Back {generate_random_string(5)}",
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
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Update card failed: {response.status_code}")
    @task(1)
    def delete_card(self):
        decks_with_cards = {deck_id: cards for deck_id, cards in self.decks.items() if cards}
        if not decks_with_cards:
            return

        deck_id = random.choice(list(decks_with_cards.keys()))
        card_index = random.randrange(len(self.decks[deck_id]))
        card_id = self.decks[deck_id][card_index]  # Get but don't pop yet
        
        headers = {"Authorization": f"Bearer {self.auth_token}"} if self.auth_token else {}
        with self.client.delete(
            f"/api/v1/decks/{deck_id}/cards/{card_id}/",
            headers=headers,
            catch_response=True,
            name="DELETE Card"
        ) as response:
            if response.status_code == 204:
                # Only remove from tracking if delete was successful
                self.decks[deck_id].pop(card_index)
                if card_id in self.card_ids:
                    self.card_ids.remove(card_id)
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
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
    weight = 2
    
    def on_start(self):
        self.user_data = random.choice(USERS) if USERS else None
        if not self.user_data:
            return
        self.auth_token = self.user_data['token']
    
    @task(10)
    def read_health(self):
        with self.client.get(
            "/api/v1/status/",
            catch_response=True,
            name="ReadOnly: Health Check"
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Failed with status {response.status_code}")
    
    @task(5)
    def read_user_decks(self):
        if not self.user_data:
            return
        headers = {"Authorization": f"Bearer {self.auth_token}"}
        with self.client.get(
            "/api/v1/users/decks/",
            headers=headers,
            catch_response=True,
            name="ReadOnly: Get Decks"
        ) as response:
            if response.status_code == 200:
                response.success()
            elif response.status_code == 429:
                response.success()  # Rate limited, but don't fail
            else:
                response.failure(f"Failed with status {response.status_code}")

class WriteHeavyUser(HttpUser):
    """User that performs write-heavy operations"""
    wait_time = between(2, 5)  # Longer wait time to avoid rate limiting
    weight = 1
    
    def on_start(self):
        self.user_data = random.choice(USERS) if USERS else None
        if not self.user_data:
            self.stop()
            return
        self.auth_token = self.user_data['token']
    
    @task
    def create_many_decks(self):
        """Create multiple decks with proper rate limiting"""
        if not self.user_data:
            return
        
        for i in range(3):
            payload = {
                "title": f"Bulk Deck {generate_random_string(5)}"
            }
            headers = {"Authorization": f"Bearer {self.auth_token}"}
            with self.client.post(
                "/api/v1/decks",
                json=payload,
                headers=headers,
                catch_response=True,
                name="WriteHeavy: Create Deck"
            ) as response:
                if response.status_code == 201:
                    response.success()
                elif response.status_code == 429:
                    response.success()  # Rate limited, stop but don't fail
                    break
                else:
                    response.failure(f"Failed to create deck: {response.status_code} - {response.text}")