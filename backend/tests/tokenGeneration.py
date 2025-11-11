import firebase_admin
from firebase_admin import credentials, auth
import json
import os
import sys
import requests

class FirebaseTokenGenerator:
    def __init__(self, service_account_path="../test_service-account-key.json"):        
        if not firebase_admin._apps:
            cred = credentials.Certificate(service_account_path)
            firebase_admin.initialize_app(cred)
    
    def create_custom_token(self, uid, additional_claims=None):
        return auth.create_custom_token(uid, additional_claims)
    
    def exchange_custom_token_for_id_token(self, custom_token):        
        resp = requests.post(
            "https://identitytoolkit.googleapis.com/v1/accounts:signInWithCustomToken",
            params={"key": os.getenv("FIREBASE_API_KEY", "")},
            json={"token": custom_token.decode("utf-8"), "returnSecureToken": True},
        )
        if resp.status_code == 200:
            return resp.json()["idToken"]
        else:
            print("Error exchanging token:", resp.text)
            return None
    def register_user_in_backend(self, id_token):
        headers = {'Authorization': f"Bearer {id_token}"}
        payload = {"name": "Stress Test User"}
        
        resp = requests.post(
            "http://localhost:8080/api/v1/users/",
            headers=headers,
            json=payload
        )
        if resp.status_code == 201:
            print("User registered successfully.")
        else:
            print("Error registering user:", resp.text)

    def create_user_and_token(self, email=None, password=None, display_name=None):
        try:
            user = auth.create_user(
                email=email,
                password=password,
                display_name=display_name
            )
            custom_token = auth.create_custom_token(user.uid)

            id_token = self.exchange_custom_token_for_id_token(custom_token)
            if id_token:
                self.register_user_in_backend(id_token)
            return user.uid, id_token

        except Exception as e:
            print(f"Error creating user: {e}")
            return None, None
    
    def create_multiple_users(self, count=10):
        users=[]
        for i in range(count):
            email = f"stresstest_{i}@example.com"
            password = "password"
            display_name = f"Stress Test User {i}"
            uid, token = self.create_user_and_token(email, password, display_name)
            if uid:
                users.append({"uid": uid, "token": token})
        return users

    def delete_user(self, uid):
        try:
            auth.delete_user(uid)
        except Exception as e:
            print(f"Error deleting user {uid}: {e}")

    def delete_multiple_users(self, user_data):
        for user in user_data:
            self.delete_user(user['uid'])

if __name__ == "__main__":
    import argparse
    
    parser = argparse.ArgumentParser(description='Generate Firebase tokens for stress testing')
    parser.add_argument('--users', type=int, default=1, help='Number of users to create (default: 1)')
    parser.add_argument('--delete-file', type=str, help='Delete users from JSON file')
    
    args = parser.parse_args()
    
    generator = FirebaseTokenGenerator()
    
    if args.delete_file:
        try:
            with(open(args.delete_file, 'r') as f):
                user_data = json.load(f)
            generator.delete_multiple_users(user_data)
            os.remove(args.delete_file)
            print(f"Deleted users and removed {args.delete_file}")
        except Exception as e:
            print(f"Error deleting users from file: {e}")
            sys.exit(1)
        sys.exit(0)

    if args.users == 1:
        uid, token = generator.create_user_and_token(
            email="test@example.com",
            password="password",
            display_name="Test User"
        )
        
        if uid and token:
            print("Your Firebase ID Token:\n")
            print(token)
    else:
        users = generator.create_multiple_users(args.users)
        output_file = "stress_test_users.json"
        try:
            with open(output_file, 'w') as f:
                json.dump(users, f, indent=4)
            print(f"Created {args.users} users. Details saved in {output_file}")
        except Exception as e:
            print(f"Error saving users to file: {e}")
            sys.exit(1)