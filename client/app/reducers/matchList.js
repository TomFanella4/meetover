import {
  FETCH_MATCHES,
  FETCH_PROFILE
} from '../actions/actionTypes';

// TODO combine matches and profiles
const initialMatchList = {
  matches: [],
  profile: {},
};

export default matchList = (state = initialMatchList, action) => {
  switch(action.type) {
    case FETCH_MATCHES:
      return { ...state, matches: action.matches };

    case FETCH_PROFILE:
      return { ...state, profile: action.profile };

    default:
      return state;
  }
};
