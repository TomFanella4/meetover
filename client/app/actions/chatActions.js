import { serverURI } from '../constants/Common';
import {
  registerFetchFirebaseThreadList,
  registerFetchFirebaseNewMessage,
  fetchFirebaseEarlierMessages
} from '../firebase';

import {
  FETCH_THREAD_LIST,
  FETCH_NEW_MESSAGE,
  FETCH_EARLIER_MESSAGES
} from './actionTypes';

const useMocks = true;

const fetchThreadList = threadList => ({
  type: FETCH_THREAD_LIST,
  threadList
});

const fetchNewMessage = (_id, message) => ({
  type: FETCH_NEW_MESSAGE,
  _id,
  message
});

const fetchEarlierMessages = (_id, messages) => ({
  type: FETCH_EARLIER_MESSAGES,
  _id,
  messages
});

export const registerFetchThreadListAsync = () => {
  return async (dispatch, getState) => {
    registerFetchFirebaseThreadList(async threadList => {
      const state = getState();
      const accessToken = state.userProfile.token.access_token;
      const userId = state.userProfile.id;
      const ids = [];
      const threads = Object.values(threadList);

      threads.forEach(thread => !thread.profile && ids.push(thread.userID));
      const uri = `${serverURI}/userprofiles`;
      const init = {
        method: 'POST',
        body: JSON.stringify(ids),
        headers: new Headers({
          'Token': accessToken,
          'Identity': userId
        })
      };

      const response = await fetch(uri, init)
        .catch(err => console.log(err));
      const result = await response.json()
        .catch(err => console.log(err));

      threads.forEach(thread => {
        if (result[thread.userID]) {
          threadList[thread._id].profile = result[thread.userID].profile;
        }
      });

      dispatch(fetchThreadList(threadList));
    });
  };
};

export const registerFetchNewMessageAsync = _id => {
  return async dispatch => {
    registerFetchFirebaseNewMessage(
      _id,
      message => dispatch(fetchNewMessage(_id, message))
    );
  };
};

export const fetchEarlierMessagesAsync = (_id, limit, endId) => {
  return async dispatch => {
    fetchFirebaseEarlierMessages(
      _id,
      messages => dispatch(fetchEarlierMessages(_id, messages)),
      limit,
      endId
    );
  };
};
