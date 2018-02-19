import { combineReducers } from 'redux';

import matchList from './matchList'
import navigation from './navigation';
import userProfile from './userProfile';

const rootReducer = combineReducers({
  matchList,
  navigation,
  userProfile
  // Add future reducers here
});

export default rootReducer;
