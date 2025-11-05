"use client"
import { createContext, useContext, useEffect, useState } from "react"

import { initializeApp } from "firebase/app";
import { AuthProvider as AuthProviderType, getAuth, signInWithPopup, User } from "firebase/auth";
import { GoogleAuthProvider } from "firebase/auth";


export const providers = {
  google: new GoogleAuthProvider(),
};

// TODO: Replace the following with your app's Firebase project configuration
// See: https://firebase.google.com/docs/web/learn-more#config-object
const firebaseConfig = {
  apiKey: "AIzaSyB3YgTG_jIt-d2hT03c3GKc-LucVjkBLlY",
  authDomain: "integrasjonsprosjekt-c5a18.firebaseapp.com",
  projectId: "integrasjonsprosjekt-c5a18",
  storageBucket: "integrasjonsprosjekt-c5a18.firebasestorage.app",
  messagingSenderId: "28577550135",
  appId: "1:28577550135:web:63f269e47067b7d72814b2",
  measurementId: "G-KMXWVG89QB"
};

// Initialize Firebase
const app = initializeApp(firebaseConfig);


// Initialize Firebase Authentication and get a reference to the service
const auth = getAuth(app);

interface AuthContextType {
  user: User | null;
  signIn: (provider: AuthProviderType) => Promise<void>;
  signOut: () => Promise<void>;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  
  async function signIn(provider: AuthProviderType) {
    const result = await signInWithPopup(auth, provider);
    setUser(result.user);
  }

  async function signOut() {
    await auth.signOut();
    setUser(null);
  }

  auth.authStateReady().then(() => {
    setUser(auth.currentUser);
  });

  return (
    <AuthContext value={{ user, signIn, signOut }}>
      {children}
    </AuthContext>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
