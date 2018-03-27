import { combineReducers } from 'redux';

import matchList from './matchList'
import userProfile from './userProfile';

const rootReducer = combineReducers({
  matchList,
  userProfile
});

export default rootReducer;
