'server-only';
import { cache } from 'react';
import { FirebaseApp, getApp, getApps, initializeApp, initializeServerApp } from 'firebase/app';
import { getAuth, User } from 'firebase/auth';
import { firebaseConfig } from './app';
import { cookies } from 'next/headers';

function getOrInitApp(): FirebaseApp {
  const apps = getApps();
  return apps.length ? getApp() : initializeApp(firebaseConfig);
}

export const getUser = cache(async (): Promise<User | null> => {
  const idToken = (await cookies()).get('__session')?.value;
  const baseApp = getOrInitApp();
  const serverApp = initializeServerApp(firebaseConfig, { authIdToken: idToken });
  const auth = getAuth(serverApp);

  await auth.authStateReady();

  const u = auth.currentUser;
  return u;
});
