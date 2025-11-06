import { getAuth, GithubAuthProvider, GoogleAuthProvider, signInWithPopup } from "firebase/auth";
import { app } from "@/lib/firebase/app";
import { deleteCookie, setCookie } from "cookies-next";

// Initialize Firebase Authentication and get a reference to the service
export const auth = getAuth(app);

auth.onIdTokenChanged(async user => {
  if (user) {
    const idToken = await user.getIdToken();
    setCookie("__session", idToken);
  }
  else {
    deleteCookie("__session");
  }
})

export const providers = {
  google: new GoogleAuthProvider(),
  github: new GithubAuthProvider(),
};

export async function signIn(provider: keyof typeof providers) {
  try {
    const result = await signInWithPopup(auth, providers[provider]);
    return result.user;
  } catch (error) {
    console.error("Error signing in:", error);
    throw error;
  }
}

export async function signOut() {
  await auth.signOut();
}

export const getCurrentUser = () => {
  return auth.currentUser;
}
