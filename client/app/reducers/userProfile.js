import {
  LOGIN,
  MODIFY_USER_PROFILE,
  LOGOUT
} from '../actions/actionTypes';

const initalUserProfile = {
  isAuthenticated: false,
};

export default userProfile = (state = initalUserProfile, action) => {
  switch (action.type) {

    case LOGIN:
      return { isAuthenticated: true, ...action.userProfile };

    case MODIFY_USER_PROFILE:
      return { ...state, ...action.userProfile };

    case LOGOUT:
      return initalUserProfile;

    default:
      return state;
  }
};
