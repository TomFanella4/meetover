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
  return async dispatch => {
    registerFetchFirebaseThreadList(
      threadList => dispatch(fetchThreadList(threadList))
    );
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
