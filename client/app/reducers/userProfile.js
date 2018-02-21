import {
  CREATE_PROFILE,
  LOGIN,
  MODIFY_USER_PROFILE,
  LOGOUT
} from '../actions/actionTypes';

const initalUserProfile = {
  isAuthenticated: false,
  shareLocation: true
};

export default userProfile = (state = initalUserProfile, action) => {
  switch (action.type) {

    case CREATE_PROFILE:
    case LOGIN:
    case MODIFY_USER_PROFILE:
      return { ...state, ...action.userProfile };

    case LOGOUT:
      return initalUserProfile;

    default:
      return state;
  }
};
