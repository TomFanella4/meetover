import Expo from 'expo';

import {
  FETCH_MATCHES,
  LOGIN,
  LOGOUT,
  UPDATE_USER_PROFILE,
} from './actionTypes';

import matchesMock from '../mocks/matches';

const useMocks = true;

export const fetchMatches = matches => ({
  type: FETCH_MATCHES,
  matches
});

export const fetchMatchesAsync = userId => {
  return async dispatch => {
    let matches;

    if (useMocks) {
      matches = matchesMock;
    } else {
      const uri = `https://meetover.herokuapp.com/match/${userId}`;
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      matches = await response.json();
    }

    dispatch(fetchMatches(matches));
  };
};

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
