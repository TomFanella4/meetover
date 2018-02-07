import { combineReducers } from 'redux';

import userProfileReducer from './userProfileReducer';

const rootReducer = combineReducers({
  userProfileReducer
  // Add future reducers here
});

export default rootReducer;
