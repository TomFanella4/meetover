import Expo, { Permissions } from 'expo';
import { times } from 'lodash';
import { serverURI } from '../constants/Common';

import {
  FETCH_MATCHES
} from './actionTypes';

const useMocks = false;

const fetchMatches = matches => ({
  type: FETCH_MATCHES,
  matches
});

export const fetchMatchesAsync = (userId, accessToken) => {
  return async dispatch => {

    const uri = `${serverURI}/match/${userId}`;
    const init = {
      method: 'POST',
      headers: new Headers({
        'Token': accessToken,
        'Identity': userId
      })
    };

    const response = await fetch(uri, init)
      .catch(err => console.log(err));
    const { matches } = await response.json()
      .catch(err => console.log(err));

    dispatch(fetchMatches(matches));
  };
};
