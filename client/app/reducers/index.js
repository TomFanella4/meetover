import { combineReducers } from 'redux';

import matchList from './matchList';
import chat from './chat';
import userProfile from './userProfile';

const rootReducer = combineReducers({
  matchList,
  chat,
  userProfile
});

export default rootReducer;
