import * as firebase from 'firebase/app';
import 'firebase/auth';
import 'firebase/database';

import { FIREBASE_API_KEY, FIREBASE_SENDER_ID } from 'react-native-dotenv';

const config = {
  apiKey: FIREBASE_API_KEY,
  authDomain: "meetoverdb.firebaseapp.com",
  databaseURL: "https://meetoverdb.firebaseio.com",
  projectId: "meetoverdb",
  storageBucket: "meetoverdb.appspot.com",
  messagingSenderId: FIREBASE_SENDER_ID
};

firebase.initializeApp(config);

export const signInToFirebase = async (token, accessToken) => {
  try {
    return await firebase.auth().signInWithCustomToken(token);
  } catch (error) {
    const uri = `https://meetover.herokuapp.com/login/refresh/${accessToken}`
    const init = { method: 'POST' };

    const response = await fetch(uri, init);
    const { firebaseCustomToken } = await response.json();

    try {
      return await firebase.auth().signInWithCustomToken(firebaseCustomToken);
    } catch (err) {
      console.log(`Could not sign in to Firebase: ${err}`);
    }
  }
};

export async function fetchIdToken(token, accessToken){
  let user = firebase.auth().currentUser;
  if (!user) {
    await signInToFirebase(token, accessToken);
    user = firebase.auth().currentUser;
  }

  return await user.getIdToken(true)
    .catch(err => {
      console.log(`Could not fetch Firebase ID Token: ${err}`);

      throw err;
    });
};

export const modifyFirebaseUserProfile = async (key, value) => {
  const user = firebase.auth().currentUser;
  if (!user) {
    throw 'User is not signed in';
  }
  return await firebase.database().ref(`users/${user.uid}/profile`).update({
    [key]: value
  });
};
