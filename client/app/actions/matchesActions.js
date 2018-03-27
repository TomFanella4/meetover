import { times } from 'lodash';

import {
  FETCH_MATCHES,
  FETCH_PROFILE,
} from './actionTypes';

const useMocks = true;

const fetchMatches = matches => ({
  type: FETCH_MATCHES,
  matches
});

const fetchProfile = profile => ({
  type: FETCH_PROFILE,
  profile
});

export const fetchMatchesAsync = userId => {
  return async dispatch => {
    let matches;

    if (useMocks) {
      const uri = 'https://meetover.herokuapp.com/test/liprofile';
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      let profile = await response.json();
      profile = JSON.parse(profile); // TODO make all the JSON.parse unnecessary

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

export const fetchProfileAsync = userId => {
  return async dispatch => {
    let profile;

    if (useMocks) {
      const uri = 'https://meetover.herokuapp.com/test/liprofile';
      const init = { method: 'POST' };

      const response = await fetch(uri, init);
      profile = await response.json();
      profile = JSON.parse(profile);
    } else {
      const uri = `https://meetover.herokuapp.com/people/${userId}`;
      const init = { method: 'GET' };

      const response = await fetch(uri, init);
      profile = await response.json();
      profile = JSON.parse(profile);
    }

    dispatch(fetchProfile(profile));
  };
};
