import * as firebase from 'firebase/app';
import 'firebase/auth';

import { FIREBASE_API_KEY, FIREBASE_SENDER_ID } from 'react-native-dotenv';

const config = {
  apiKey: FIREBASE_API_KEY,
  authDomain: "meetoverdb.firebaseapp.com",
  databaseURL: "https://meetoverdb.firebaseio.com",
  projectId: "meetoverdb",
  storageBucket: "meetoverdb.appspot.com",
  messagingSenderId: FIREBASE_SENDER_ID
};

const app = firebase.initializeApp(config);

export async function fetchIdToken(token){
  await firebase.auth().signInWithCustomToken(token)
    .catch(err => {
      console.log(`Could not sign in to Firebase: ${err}`);

      throw err;
    });

  return await firebase.auth().currentUser.getIdToken(true)
    .catch(err => {
      console.log(`Could not fetch Firebase ID Token: ${err}`);

      throw err;
    });
};
