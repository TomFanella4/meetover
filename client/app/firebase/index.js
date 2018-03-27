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

export const signInToFirebase = async (token, accessToken, id) => {
  try {
    await firebase.auth().signInWithCustomToken(token);
    return token;
  } catch (error) {
    const uri = `https://meetover.herokuapp.com/login/refresh/${accessToken}`
    const init = {
      method: 'POST',
      headers: new Headers({
        'Token': accessToken,
        'Identity': id
      })
    };

    const response = await fetch(uri, init);

    if (response.status === 401) {
      const err = 'Could not sign in to Firebase: invalid credentials';
      console.log(err);
      throw err;
    }

    const { firebaseCustomToken } = await response.json();

    try {
      await firebase.auth().signInWithCustomToken(firebaseCustomToken);
      return firebaseCustomToken;
    } catch (err) {
      console.log(`Could not sign in to Firebase: ${err}`);
      throw err;
    }
  }
};

export async function fetchIdToken(token) {
  return await firebase.auth().currentUser.getIdToken(true)
    .catch(err => {
      console.log(`Could not fetch Firebase ID Token: ${err}`);

      throw err;
    });
};

export const modifyFirebaseUserProfile = async (key, value) => {
  const user = firebase.auth().currentUser;

  return await firebase.database().ref(`users/${user.uid}/profile`)
    .update({ [key]: value })
    .catch(err => {
      console.log(`Could not modify user profile: ${err}`);

      throw err;
    });
};
