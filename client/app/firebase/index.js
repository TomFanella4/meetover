import * as firebase from 'firebase/app';
import 'firebase/auth';
import 'firebase/database';

import { FIREBASE_API_KEY, FIREBASE_SENDER_ID } from 'react-native-dotenv';

import { serverURI } from '../constants/Common';

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
    const uri = `${serverURI}/refreshtoken`;
    const init = {
      method: 'POST',
      headers: new Headers({
        'Token': accessToken,
        'Identity': id
      })
    };

    const response = await fetch(uri, init);

    if (response.status !== 200) {
      const err = 'Could not sign in to Firebase: invalid credentials';
      console.log(err);
      console.log(response);
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

export const registerFetchFirebaseThreadList = updateFn => {
  const user = firebase.auth().currentUser;
  const threadListRef = firebase.database().ref(`users/${user.uid}/threadList`);
  threadListRef.on('value', snapshot => {
    let threadList = snapshot.val();
    threadList = threadList ? Object.values(threadList) : [];
    updateFn(threadList);
  });
};

export const registerFetchFirebaseNewMessage = (_id, updateFn) => {
  const messageRef = firebase.database().ref(`messages/${_id}`).limitToLast(1);
  messageRef.on('child_added', message => updateFn(message.val()));
};

export const fetchFirebaseEarlierMessages = (_id, updateFn, limit, endId) => {
  let messagesRef;
  endId ?
    messagesRef = firebase.database().ref(`messages/${_id}`)
    .endAt(null, endId).limitToLast(limit + 1)
  :
    messagesRef = firebase.database().ref(`messages/${_id}`).limitToLast(limit);

  messagesRef.once('value', snapshot => {
    const messages = Object.values(snapshot.val()).reverse();
    endId && messages.shift();
    updateFn(messages);
  });
};

export const sendFirebaseMessage = (_id, messages) => {
  const messagesRef = firebase.database().ref(`messages/${_id}`);

  messages.forEach(message => {
    const messageRef = messagesRef.push();
    message.createdAt = message.createdAt.toISOString();
    message._id = messageRef.key;
    messageRef.set(message);
  });
};

export const chatThreadExists = async (threadId) => {
  const threadRef = firebase.database().ref(`messages/${threadId}`);
  let exists;

  await threadRef.once('value', snapshot => {exists = (snapshot.val() !== null)});

  return exists;
}
