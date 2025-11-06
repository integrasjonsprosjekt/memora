import { getApp, getApps, initializeApp } from 'firebase/app';

export const firebaseConfig = {
  apiKey: 'AIzaSyB3YgTG_jIt-d2hT03c3GKc-LucVjkBLlY',
  authDomain: 'integrasjonsprosjekt-c5a18.firebaseapp.com',
  projectId: 'integrasjonsprosjekt-c5a18',
  storageBucket: 'integrasjonsprosjekt-c5a18.firebasestorage.app',
  messagingSenderId: '28577550135',
  appId: '1:28577550135:web:63f269e47067b7d72814b2',
  measurementId: 'G-KMXWVG89QB',
};

// Initialize Firebase
export const app = getApps().length > 0 ? getApp() : initializeApp(firebaseConfig);
