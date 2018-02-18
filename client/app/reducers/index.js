import { combineReducers } from 'redux';

import navigation from './navigation';
import userProfile from './userProfile';

const rootReducer = combineReducers({
  navigation,
  userProfile
  // Add future reducers here
});

export default rootReducer;
