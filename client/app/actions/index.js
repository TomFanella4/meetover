import Expo from 'expo';

import {
  LOGIN,
  LOGOUT,
  UPDATE_USER_PROFILE,
} from './actionTypes';

export const login = userProfile => ({
  type: LOGIN,
  userProfile
});

export const logout = () => ({
  type: LOGOUT
});

export const modifyUserProfile = userProfile => ({
  type: UPDATE_USER_PROFILE,
  userProfile
});

export const saveProfileAndLoginAsync = userProfile => (
  dispatch => (
    Expo.SecureStore.setItemAsync('userProfile', JSON.stringify(userProfile))
    .then(() => dispatch(login(userProfile)))
    .catch(err => console.error(err))
  )
);
