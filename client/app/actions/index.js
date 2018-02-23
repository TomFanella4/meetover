import Expo from 'expo';
import { times } from 'lodash';

import {
  FETCH_MATCHES,
  CREATE_PROFILE,
  LOGIN,
  LOGOUT,
  MODIFY_USER_PROFILE,
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
      const uri = `https://meetover.herokuapp.com/test/liprofile`;
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      let profile = await response.json();
      profile = JSON.parse(profile);

      matches = times(10, () => Object.assign({}, profile));
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

export const createProfile = userProfile => ({
  type: CREATE_PROFILE,
  userProfile
});

export const modifyUserProfile = userProfile => ({
  type: MODIFY_USER_PROFILE,
  userProfile
});

export const saveProfileAndLoginAsync = userProfile => (
  dispatch => (
    Expo.SecureStore.setItemAsync('userProfile', JSON.stringify(userProfile))
    .then(() => dispatch(login(userProfile)))
    .catch(err => {
      dispatch(login(userProfile));
      // TODO Notify user of error
      console.log(err);
    })
  )
);
