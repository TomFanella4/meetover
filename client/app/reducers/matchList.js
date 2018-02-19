import {
  FETCH_MATCHES
} from '../actions/actionTypes';

const initialMatchList = {
  matches: []
};

export default matchList = (state = initialMatchList, action) => {
  switch(action.type) {
    case FETCH_MATCHES:
      return { ...state, matches: action.matches };

    default:
      return state;
  }
};
