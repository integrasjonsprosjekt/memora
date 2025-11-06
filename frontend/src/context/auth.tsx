"use client"
import { createContext, useContext, useState } from "react"
import { AuthProvider as AuthProviderType, signInWithPopup, User } from "firebase/auth";
import { setCookie, deleteCookie } from "cookies-next"
import { auth } from "@/lib/firebase/auth";

interface AuthContextType {
  user: User | null;
  setUser: (user: User | null) => void;
}

export const AuthContext = createContext<AuthContextType | undefined>(undefined)

// A initial user could be supplied from the server but it would not be the same
// because user is non-serializable and can't be transferd across the RSC boundary.
export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(auth.currentUser);

  auth.onAuthStateChanged((user) => {
    setUser(user);
  });

  return (
    <AuthContext.Provider value={{ user, setUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
}
