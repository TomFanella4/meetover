import {
  AUTHENTICATE_USER
} from '../actions/actionTypes';

const initalUserProfileState = {
  authenticated: false
};

export default userProfileReducer = (state = initalUserProfileState, action) => {
  switch (action.type) {
    case AUTHENTICATE_USER:
      return { ...state, authenticated: true };

    default:
      return state;
  }
};
