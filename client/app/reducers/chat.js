import { GiftedChat } from 'react-native-gifted-chat';
import {
  FETCH_THREAD_LIST,
  FETCH_NEW_MESSAGE,
  FETCH_EARLIER_MESSAGES
} from '../actions/actionTypes';

const initialState = {
  threadList: null,
  messageThreads: {}
};

export default chat = (state = initialState, action) => {
  let messages;
  switch(action.type) {

    case FETCH_THREAD_LIST:
      return { ...state, threadList: action.threadList };

    case FETCH_NEW_MESSAGE:
      messages = state.messageThreads[action._id];

      if (!messages || messages[0]._id === action.message._id) {
        return state;
      }

      messages = GiftedChat.append(
        messages,
        action.message
      );

      return {
        ...state,
        messageThreads: {
          ...state.messageThreads,
          [action._id]: messages
        }
      };

    case FETCH_EARLIER_MESSAGES:
      messages = GiftedChat.prepend(
        state.messageThreads[action._id],
        action.messages
      );

      return {
        ...state,
        messageThreads: {
          ...state.messageThreads,
          [action._id]: messages
        }
      };

    default:
      return state;
  }
};
